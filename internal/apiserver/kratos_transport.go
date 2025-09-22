package apiserver

import (
	"context"

	handlergrpc "github.com/ArthurWang23/miniblog/internal/apiserver/handler/grpc"
	mwgin "github.com/ArthurWang23/miniblog/internal/pkg/middleware/gin"
	mwrkrt "github.com/ArthurWang23/miniblog/internal/pkg/middleware/krt"
	apiv1 "github.com/ArthurWang23/miniblog/pkg/api/apiserver/v1"
	kratoslog "github.com/ArthurWang23/miniblog/pkg/log"
	genericvalidation "github.com/ArthurWang23/miniblog/pkg/validation"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/transport"
	kratosgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// 构建 Kratos HTTP Server，并挂载 gin.Engine
func (c *ServerConfig) NewKratosHTTPServer() transport.Server {
	engin := gin.New()
	engin.Use(gin.Recovery(), mwgin.NoCache, mwgin.Cors, mwgin.Secure, mwgin.RequestIDMiddleware())
	c.InstallRESTAPI(engin)

	// 业务无关路由
	pprof.Register(engin)
	engin.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(404, "Page not found.")
	})

	hs := kratoshttp.NewServer(
		kratoshttp.Network(c.cfg.HTTPOptions.Network),
		kratoshttp.Address(c.cfg.HTTPOptions.Addr),
		kratoshttp.Timeout(c.cfg.HTTPOptions.Timeout),
		kratoshttp.Logger(kratoslog.Kratos()),
		kratoshttp.TLSConfig(c.cfg.TLSOptions.MustTLSConfig()),
	)
	hs.HandlePrefix("/", engin)
	return hs
}

// 构建 Kratos gRPC Server，复用现有拦截器链与服务注册
func (c *ServerConfig) NewKratosGRPCServer() transport.Server {
	// Kratos middleware 链 + 白名单
	whitelist := map[string]struct{}{
		apiv1.MiniBlog_Healthz_FullMethodName:    {},
		apiv1.MiniBlog_CreateUser_FullMethodName: {},
		apiv1.MiniBlog_Login_FullMethodName:      {},
	}

	gs := kratosgrpc.NewServer(
		kratosgrpc.Network(c.cfg.GRPCOptions.Network),
		kratosgrpc.Address(c.cfg.GRPCOptions.Addr),
		kratosgrpc.Timeout(c.cfg.GRPCOptions.Timeout),
		kratosgrpc.Logger(kratoslog.Kratos()),
		kratosgrpc.TLSConfig(c.cfg.TLSOptions.MustTLSConfig()),
		// 使用 Kratos middleware 链（取代原 UnaryInterceptor）
		kratosgrpc.Middleware(
			mwrkrt.RequestID(),
			mwrkrt.Authn(c.retriever, whitelist),
			mwrkrt.Authz(c.authz, whitelist),
			mwrkrt.Defaulter(),
			mwrkrt.Validator(genericvalidation.NewValidator(c.val)),
		),
	)

	apiv1.RegisterMiniBlogServer(gs, handlergrpc.NewHandler(c.biz))
	return gs
}

// 构建 Kratos HTTP Server 作为 grpc-gateway
func (c *ServerConfig) NewKratosGatewayHTTPServer() (transport.Server, error) {
	mux := runtime.NewServeMux()

	// gRPC 客户端连接到本地 gRPC 服务
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

	if err := apiv1.RegisterMiniBlogHandler(context.Background(), mux, conn); err != nil {
		return nil, err
	}

	hs := kratoshttp.NewServer(
		kratoshttp.Network(c.cfg.HTTPOptions.Network),
		kratoshttp.Address(c.cfg.HTTPOptions.Addr),
		kratoshttp.Timeout(c.cfg.HTTPOptions.Timeout),
		kratoshttp.Logger(kratoslog.Kratos()),
		kratoshttp.TLSConfig(c.cfg.TLSOptions.MustTLSConfig()),
	)
	hs.HandlePrefix("/", mux)
	return hs, nil
}

// 根据 serverMode 返回一个或多个 Kratos Server
func (c *ServerConfig) NewKratosServers() ([]transport.Server, error) {
	switch c.cfg.ServerMode {
	case GinServerMode:
		return []transport.Server{c.NewKratosHTTPServer()}, nil
	case GRPCServerMode:
		return []transport.Server{c.NewKratosGRPCServer()}, nil
	case GRPCGatewayServerMode:
		grpcSrv := c.NewKratosGRPCServer()
		httpSrv, err := c.NewKratosGatewayHTTPServer()
		if err != nil {
			return nil, err
		}
		return []transport.Server{grpcSrv, httpSrv}, nil
	default:
		// 默认走 gRPC
		return []transport.Server{c.NewKratosGRPCServer()}, nil
	}
}
