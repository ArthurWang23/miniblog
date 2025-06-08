package apiserver

import (
	"context"

	handler "github.com/ArthurWang23/miniblog/internal/apiserver/handler/grpc"
	"github.com/ArthurWang23/miniblog/internal/pkg/server"
	apiv1 "github.com/ArthurWang23/miniblog/pkg/api/apiserver/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type grpcServer struct {
	srv server.Server

	stop func(context.Context)
}

var _ server.Server = (*grpcServer)(nil)

// 创建并初始化grpc或grpc和grpc-gateway服务器
// NewGRPCServerOr中Or一般表示或者
// 暗示函数会有两种或多种选择中选择一种可能性

func (c *ServerConfig) NewGRPCServerOr() (server.Server, error) {
	grpcsrv, err := server.NewGRPCServer(
		c.cfg.GRPCOptions,
		func(s grpc.ServiceRegistrar) {
			apiv1.RegisterMiniBlogServer(s, handler.NewHandler())
		},
	)
	if err != nil {
		return nil, err
	}
	if c.cfg.ServerMode == GRPCServerMode {
		return &grpcServer{
			srv: grpcsrv,
			stop: func(ctx context.Context) {
				grpcsrv.GracefulStop(ctx)
			},
		}, nil
	}
	// grpc + grpc-gateway
	// 先启动grpc服务器，因为http服务器依赖grpc服务器
	go grpcsrv.RunOrDie()

	httpsrv, err := server.NewGRPCGatewayServer(
		c.cfg.HTTPOptions,
		c.cfg.GRPCOptions,
		func(mux *runtime.ServeMux, conn *grpc.ClientConn) error {
			return apiv1.RegisterMiniBlogHandler(context.Background(), mux, conn)
		},
	)
	if err != nil {
		return nil, err
	}

	return &grpcServer{
		srv: httpsrv,
		stop: func(ctx context.Context) {
			grpcsrv.GracefulStop(ctx)
			httpsrv.GracefulStop(ctx)
		},
	}, nil
}

func (s *grpcServer) RunOrDie() {
	s.srv.RunOrDie()
}

func (s *grpcServer) GracefulStop(ctx context.Context) {
	s.stop(ctx)
}
