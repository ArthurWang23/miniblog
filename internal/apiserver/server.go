// Copyright 2025 ArthurWang &lt;2826979176@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/arthurwang23/miniblog. The professional
// version of this repository is https://github.com/arthurwang23/miniblog.

package apiserver

// 运行时代码保存在internal/apiserver
// 初始化配置正确加载并读取后，基于初始化配置创建运行时配置，并基于运行时配置创建服务器实例
// 采用面向对象风格UnionServer结构体封装服务相关功能

import (
	"context"
	"github.com/glebarez/sqlite"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ArthurWang23/miniblog/internal/apiserver/biz"
	"github.com/ArthurWang23/miniblog/internal/apiserver/model"
	"github.com/ArthurWang23/miniblog/internal/apiserver/pkg/validation"
	"github.com/ArthurWang23/miniblog/internal/apiserver/store"
	"github.com/ArthurWang23/miniblog/internal/pkg/contextx"
	"github.com/ArthurWang23/miniblog/internal/pkg/known"
	"github.com/ArthurWang23/miniblog/internal/pkg/log"
	mw "github.com/ArthurWang23/miniblog/internal/pkg/middleware/gin"
	"github.com/ArthurWang23/miniblog/internal/pkg/server"
	"github.com/ArthurWang23/miniblog/pkg/auth"
	genericoptions "github.com/ArthurWang23/miniblog/pkg/options"
	"github.com/ArthurWang23/miniblog/pkg/store/where"
	"github.com/ArthurWang23/miniblog/pkg/token"
	"github.com/onexstack/onexstack/pkg/ptr"
	"gorm.io/gorm"
)

const (
	// GRPCServerMode 定义 gRPC 服务模式.
	// 使用 gRPC 框架启动一个 gRPC 服务器.
	GRPCServerMode = "grpc"
	// GRPCServerMode 定义 gRPC + HTTP 服务模式.
	// 使用 gRPC 框架启动一个 gRPC 服务器 + HTTP 反向代理服务器.
	GRPCGatewayServerMode = "grpc-gateway"
	// GinServerMode 定义 Gin 服务模式.
	// 使用 Gin Web 框架启动一个 HTTP 服务器.
	GinServerMode = "gin"
)

type Config struct {
	ServerMode        string
	JWTKey            string
	Expiration        time.Duration
	GRPCOptions       *genericoptions.GRPCOptions
	HTTPOptions       *genericoptions.HTTPOptions
	MySQLOptions      *genericoptions.MySQLOptions
	TLSOptions        *genericoptions.TLSOptions
	EnableMemoryStore bool
}

// 根据ServerMode决定要启动的服务器类型
type UnionServer struct {
	srv server.Server
}

type ServerConfig struct {
	cfg *Config
	biz biz.IBiz
	val *validation.Validator

	retriever mw.UserRetriever
	authz     *auth.Authz
}

func (cfg *Config) NewUnionServer() (*UnionServer, error) {
	// 注册租户解析函数，通过上下文获取用户id
	where.RegisterTenant("userID", func(ctx context.Context) string {
		return contextx.UserID(ctx)
	})

	// 初始化token
	token.Init(cfg.JWTKey, known.XUserID, cfg.Expiration)

	log.Infow("Initializing federation server", "server-mode", cfg.ServerMode)

	// 创建服务配置，这些配置可用来创建服务器
	srv, err := InitializeWebServer(cfg)
	if err != nil {
		return nil, err
	}
	return &UnionServer{srv: srv}, nil
}

// 启动服务并优雅关闭
func (s *UnionServer) Run() error {
	go s.srv.RunOrDie()

	// 执行kill时默认发送SIGTERM
	// 使用kill -2 发送SIGINT（如Ctrl+C）
	// 使用kill -9 发送SIGKILL，但该信号无法被捕获，因此无需监听和处理
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// 阻塞 等待从quit channel接收到信号
	<-quit

	log.Infow("Shutting down sever ...")

	// 优雅关闭服务
	// 创建上下文对象ctx，为优雅关闭服务提供超时控制
	// 确保服务在一定时间内完成清理工作
	// 若超时指定时间，服务将强制终止
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 先关闭依赖的服务，再关闭被依赖的服务
	s.srv.GracefulStop(ctx)

	log.Infow("Server exited")
	return nil
}

func (cfg *Config) NewDB() (*gorm.DB, error) {
	// 可以选择是否选择内存数据库
	if !cfg.EnableMemoryStore {
		log.Infow("Initializing database connection", "type", "mysql", "addr", cfg.MySQLOptions.Addr)
		return cfg.MySQLOptions.NewDB()
	}
	log.Infow("Initializing database connection", "type", "memory", "engine", "SQLite")
	// 使用SQLite内存模式配置数据库
	// ?cache=shared 用于设置 SQLite 的缓存模式为 共享缓存模式 (shared)。
	// 默认情况下，SQLite 的每个数据库连接拥有自己的独立缓存，这种模式称为 专用缓存 (private)。
	// 使用 共享缓存模式 (shared) 后，不同连接可以共享同一个内存中的数据库和缓存。
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Errorw("Failed to create database connection", "err", err)
		return nil, err
	}
	// 自动迁移数据库结构
	if err := db.AutoMigrate(&model.UserM{}, &model.PostM{}, &model.CasbinRuleM{}); err != nil {
		log.Errorw("Failed to migrate database schema", "err", err)
		return nil, err
	}
	// 注意：这里仅仅为了实现快速部署，降低学习难度。
	// 在真实企业开发中，不能再代码中硬编码这些初始化配置，
	// 尤其是硬编码密码、密钥之类的信息.
	// 插入 casbin_rule 表记录
	adminR, userR := "role::admin", "role::user"
	casbinRules := []model.CasbinRuleM{
		{PType: ptr.To("g"), V0: ptr.To("user-000000"), V1: &adminR},
		{PType: ptr.To("p"), V0: &adminR, V1: ptr.To("*"), V2: ptr.To("*"), V3: ptr.To("allow")},
		{PType: ptr.To("p"), V0: &userR, V1: ptr.To("/v1.MiniBlog/DeleteUser"), V2: ptr.To("CALL"), V3: ptr.To("deny")},
		{PType: ptr.To("p"), V0: &userR, V1: ptr.To("/v1.MiniBlog/ListUser"), V2: ptr.To("CALL"), V3: ptr.To("deny")},
		{PType: ptr.To("p"), V0: &userR, V1: ptr.To("/v1/users"), V2: ptr.To("GET"), V3: ptr.To("deny")},
		{PType: ptr.To("p"), V0: &userR, V1: ptr.To("/v1/users/*"), V2: ptr.To("DELETE"), V3: ptr.To("deny")},
	}

	if err := db.Create(&casbinRules).Error; err != nil {
		log.Fatalw("Failed to insert casbin_rule records", "err", err)
		return nil, err
	}

	// 插入默认用户（root用户）
	user := model.UserM{
		UserID:    "user-000000",
		Username:  "root",
		Password:  "miniblog1234",
		Nickname:  "administrator",
		Email:     "colin404@foxmail.com",
		Phone:     "18110000000",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Create(&user).Error; err != nil {
		log.Fatalw("Failed to insert default root user", "err", err)
		return nil, err
	}

	return db, nil
}

type UserRetriever struct {
	store store.IStore
}

func (r *UserRetriever) GetUser(ctx context.Context, userID string) (*model.UserM, error) {
	return r.store.User().Get(ctx, where.F("userID", userID))
}

func ProviderDB(cfg *Config) (*gorm.DB, error) {
	return cfg.NewDB()
}

func NewWebServer(serverMode string, serverConfig *ServerConfig) (server.Server, error) {
	switch serverMode {
	case GinServerMode:
		return serverConfig.NewGinServer(), nil
	default:
		return serverConfig.NewGRPCServerOr()
	}
}
