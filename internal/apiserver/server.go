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
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ArthurWang23/miniblog/internal/apiserver/biz"
	"github.com/ArthurWang23/miniblog/internal/apiserver/pkg/validation"
	"github.com/ArthurWang23/miniblog/internal/apiserver/store"
	"github.com/ArthurWang23/miniblog/internal/pkg/contextx"
	"github.com/ArthurWang23/miniblog/internal/pkg/log"
	"github.com/ArthurWang23/miniblog/internal/pkg/server"
	genericoptions "github.com/ArthurWang23/miniblog/pkg/options"
	"github.com/ArthurWang23/miniblog/pkg/store/where"
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
	ServerMode   string
	JWTKey       string
	Expiration   time.Duration
	GRPCOptions  *genericoptions.GRPCOptions
	HTTPOptions  *genericoptions.HTTPOptions
	MySQLOptions *genericoptions.MySQLOptions
}

// 根据ServerMode决定要启动的服务器类型
type UnionServer struct {
	srv server.Server
}

type ServerConfig struct {
	cfg *Config
	biz biz.IBiz
	val *validation.Validator
}

func (cfg *Config) NewUnionServer() (*UnionServer, error) {
	// 注册租户解析函数，通过上下文获取用户id
	where.RegisterTenant("userID", func(ctx context.Context) string {
		return contextx.UserID(ctx)
	})
	serverConfig, err := cfg.NewServerConfig()
	if err != nil {
		return nil, err
	}
	log.Infow("Initializing federation server", "server-mode", cfg.ServerMode)

	// 根据服务模式创建对应的服务实例
	var srv server.Server
	switch cfg.ServerMode {
	case GinServerMode:
		srv, err = serverConfig.NewGinServer(), nil
	default:
		srv, err = serverConfig.NewGRPCServerOr()
	}
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

func (cfg *Config) NewServerConfig() (*ServerConfig, error) {
	db, err := cfg.NewDB()
	if err != nil {
		return nil, err
	}
	store := store.NewStore(db)
	return &ServerConfig{
		cfg: cfg,
		biz: biz.NewBiz(store),
		val: validation.New(store),
	}, nil
}

func (cfg *Config) NewDB() (*gorm.DB, error) {
	return cfg.MySQLOptions.NewDB()
}
