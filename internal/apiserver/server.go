package apiserver

// 运行时代码保存在internal/apiserver
// 初始化配置正确加载并读取后，基于初始化配置创建运行时配置，并基于运行时配置创建服务器实例
// 采用面向对象风格UnionServer结构体封装服务相关功能

import (
	"net"
	"time"

	handler "github.com/ArthurWang23/miniblog/internal/apiserver/handler/grpc"
	"github.com/ArthurWang23/miniblog/internal/pkg/log"
	apiv1 "github.com/ArthurWang23/miniblog/pkg/api/apiserver/v1"
	genericclioptions "github.com/onexstack/onexstack/pkg/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
}

type UnionServer struct {
	cfg *Config
	srv *grpc.Server
	lis net.Listener
}

func (cfg *Config) NewUnionServer() (*UnionServer, error) {
	lis, err := net.Listen("tcp", cfg.GRPCOptions.Addr)
	if err != nil {
		log.Fatalw("Failed to listen", "err", err)
		return nil, err
	}
	grpcsrv := grpc.NewServer()
	apiv1.RegisterMiniBlogServer(grpcsrv, handler.NewHandler())
	reflection.Register(grpcsrv)
	return &UnionServer{cfg: cfg, srv: grpcsrv, lis: lis}, nil
}

func (s *UnionServer) Run() error {
	log.Infow("Start to listen the incoming requests on grpc address", "addr", s.cfg.GRPCOptions.Addr)
	return s.srv.Serve(s.lis)
}
