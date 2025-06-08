package apiserver

// 运行时代码保存在internal/apiserver
// 初始化配置正确加载并读取后，基于初始化配置创建运行时配置，并基于运行时配置创建服务器实例
// 采用面向对象风格UnionServer结构体封装服务相关功能

import (
	"time"

	"github.com/ArthurWang23/miniblog/internal/pkg/log"
	"github.com/ArthurWang23/miniblog/internal/pkg/server"
	genericclioptions "github.com/onexstack/onexstack/pkg/options"
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
	ServerMode  string
	JWTKey      string
	Expiration  time.Duration
	GRPCOptions *genericclioptions.GRPCOptions
	HTTPOptions *genericclioptions.HTTPOptions
}

// 根据ServerMode决定要启动的服务器类型
type UnionServer struct {
	srv server.Server
}

type ServerConfig struct {
	cfg *Config
}

func (cfg *Config) NewUnionServer() (*UnionServer, error) {
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

func (cfg *Config) NewServerConfig() (*ServerConfig, error) {
	return &ServerConfig{cfg: cfg}, nil
}

func (s *UnionServer) Run() error {
	s.srv.RunOrDie()
	return nil
}
