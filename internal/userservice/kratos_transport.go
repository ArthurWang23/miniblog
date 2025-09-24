package userservice

import (
	mwgin "github.com/ArthurWang23/miniblog/internal/pkg/middleware/gin"
	mwrkrt "github.com/ArthurWang23/miniblog/internal/pkg/middleware/krt"
	kratoslog "github.com/ArthurWang23/miniblog/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/transport"
	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
	kratosgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	pb "github.com/ArthurWang23/miniblog/pkg/api/userservice/v1"
	handlergrpc "github.com/ArthurWang23/miniblog/internal/userservice/handler/grpc"
	handlerhttp "github.com/ArthurWang23/miniblog/internal/userservice/handler/http"
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// HTTP Server（Gin）
func (c *ServerConfig) newHTTPServer() transport.Server {
	engine := gin.New()
	engine.Use(gin.Recovery(), mwgin.NoCache, mwgin.Cors, mwgin.Secure)
	// 注入依赖到 HTTP handler
	httpHandler := handlerhttp.NewHandler(c.userStore, c.kafkaWriter)
	InstallRoutes(engine, httpHandler)

	// HTTP 路径白名单
	httpWhitelist := map[string]struct{}{
		"/healthz":  {},
		"/v1/users": {},
		"/v1/login": {},
	}

	hs := kratoshttp.NewServer(
		kratoshttp.Network(c.cfg.HTTPOptions.Network),
		kratoshttp.Address(c.cfg.HTTPOptions.Addr),
		kratoshttp.Timeout(c.cfg.HTTPOptions.Timeout),
		kratoshttp.Logger(kratoslog.Kratos()),
		kratoshttp.TLSConfig(c.cfg.TLSOptions.MustTLSConfig()),
		kratoshttp.Middleware(
			mwrkrt.RequestID(),
			// 接入认证/鉴权
			mwrkrt.Authn(c.retriever, httpWhitelist),
			mwrkrt.Authz(c.authz, httpWhitelist),
		),
	)
	hs.HandlePrefix("/", engine)
	return hs
}

// gRPC Server
func (c *ServerConfig) newGRPCServer() transport.Server {
	// gRPC FullMethod 白名单
	grpcWhitelist := map[string]struct{}{
		"/userservice.v1.UserService/Healthz":   {},
		"/userservice.v1.UserService/CreateUser": {},
		"/userservice.v1.UserService/Login":     {},
	}

	gs := kratosgrpc.NewServer(
		kratosgrpc.Network(c.cfg.GRPCOptions.Network),
		kratosgrpc.Address(c.cfg.GRPCOptions.Addr),
		kratosgrpc.Timeout(c.cfg.GRPCOptions.Timeout),
		kratosgrpc.Logger(kratoslog.Kratos()),
		kratosgrpc.TLSConfig(c.cfg.TLSOptions.MustTLSConfig()),
		kratosgrpc.Middleware(
			mwrkrt.RequestID(),
			// 接入认证/鉴权
			mwrkrt.Authn(c.retriever, grpcWhitelist),
			mwrkrt.Authz(c.authz, grpcWhitelist),
		),
	)
	// 注入依赖到 gRPC handler
	pb.RegisterUserServiceServer(gs, handlergrpc.NewHandler(c.userStore, c.kafkaWriter))
	return gs
}

// Gateway HTTP Server（grpc-gateway）
func (c *ServerConfig) newGatewayHTTPServer() (transport.Server, error) {
	mux := runtime.NewServeMux()

	var dialOpts []grpc.DialOption
	if c.cfg.TLSOptions != nil && c.cfg.TLSOptions.UseTLS {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(credentials.NewTLS(c.cfg.TLSOptions.MustTLSConfig())))
	} else {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.DialContext(context.Background(), c.cfg.GRPCOptions.Addr, dialOpts...)
	if err != nil {
		return nil, err
	}

	if err := pb.RegisterUserServiceHandler(context.Background(), mux, conn); err != nil {
		return nil, err
	}

	// HTTP 路径白名单（与 HTTP Server 保持一致）
	httpWhitelist := map[string]struct{}{
		"/healthz":  {},
		"/v1/users": {},
		"/v1/login": {},
	}

	hs := kratoshttp.NewServer(
		kratoshttp.Network(c.cfg.HTTPOptions.Network),
		kratoshttp.Address(c.cfg.HTTPOptions.Addr),
		kratoshttp.Timeout(c.cfg.HTTPOptions.Timeout),
		kratoshttp.Logger(kratoslog.Kratos()),
		kratoshttp.TLSConfig(c.cfg.TLSOptions.MustTLSConfig()),
		kratoshttp.Middleware(
			mwrkrt.RequestID(),
			// 接入认证/鉴权
			mwrkrt.Authn(c.retriever, httpWhitelist),
			mwrkrt.Authz(c.authz, httpWhitelist),
		),
	)
	hs.HandlePrefix("/", mux)
	return hs, nil
}

func (c *ServerConfig) NewKratosServers() ([]transport.Server, error) {
	switch c.cfg.ServerMode {
	case "http":
		return []transport.Server{c.newHTTPServer()}, nil
	case "grpc":
		return []transport.Server{c.newGRPCServer()}, nil
	case "both":
		return []transport.Server{c.newHTTPServer(), c.newGRPCServer()}, nil
	case "grpc-gateway":
		grpcSrv := c.newGRPCServer()
		httpSrv, err := c.newGatewayHTTPServer()
		if err != nil {
			return nil, err
		}
		return []transport.Server{grpcSrv, httpSrv}, nil
	default:
		return []transport.Server{c.newHTTPServer()}, nil
	}
}