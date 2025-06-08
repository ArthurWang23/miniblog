package apiserver

// 运行时代码保存在internal/apiserver
// 初始化配置正确加载并读取后，基于初始化配置创建运行时配置，并基于运行时配置创建服务器实例
// 采用面向对象风格UnionServer结构体封装服务相关功能

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	handler "github.com/ArthurWang23/miniblog/internal/apiserver/handler/grpc"
	"github.com/ArthurWang23/miniblog/internal/pkg/log"
	apiv1 "github.com/ArthurWang23/miniblog/pkg/api/apiserver/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	genericclioptions "github.com/onexstack/onexstack/pkg/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
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
	go s.srv.Serve(s.lis)
	// insecure.NewCredentials() 创建一个不安全的传输凭证，用于与gRPC服务器进行通信。
	// 使用这种凭据grpc客户端和服务端通信不会加密，也不会进行身份验证
	// 因为http请求转发到grpc客户端，是内部转发行为
	dialOptions := []grpc.DialOption{grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials())}

	// 创建grpc客户端连接conn
	conn, err := grpc.NewClient(s.cfg.GRPCOptions.Addr, dialOptions...)
	if err != nil {
		return err
	}
	// 创建一个http.ServeMux实例gwmux，用于处理http请求
	// UseEnumNumbers：true来设置序列化protobuf数据时，枚举类型的字段以数字格式输出，否则默认以字符串格式输出
	// gwmux grpc-gateway servemux 服务多路复用器
	// 主要是创建http到grpc的桥梁
	gwmux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseEnumNumbers: true,
		},
	}))
	if err := apiv1.RegisterMiniBlogHandler(context.Background(), gwmux, conn); err != nil {
		return err
	}
	log.Infow("Start to listen the incoming requests", "protocol", "http", "addr", s.cfg.HTTPOptions.Addr)
	httpsrv := &http.Server{
		Addr:    s.cfg.HTTPOptions.Addr,
		Handler: gwmux,
	}
	if err := httpsrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
