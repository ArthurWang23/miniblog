// Copyright 2025 ArthurWang &lt;2826979176@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/arthurwang23/miniblog. The professional
// version of this repository is https://github.com/arthurwang23/miniblog.

package apiserver

import (
	"context"

	handler "github.com/ArthurWang23/miniblog/internal/apiserver/handler/grpc"
	mw "github.com/ArthurWang23/miniblog/internal/pkg/middleware/grpc"
	"github.com/ArthurWang23/miniblog/internal/pkg/server"
	apiv1 "github.com/ArthurWang23/miniblog/pkg/api/apiserver/v1"
	genericvalidation "github.com/ArthurWang23/miniblog/pkg/validation"
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
// 添加拦截器到grpc请求链中
func (c *ServerConfig) NewGRPCServerOr() (server.Server, error) {
	// 配置拦截器链
	serverOptions := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			mw.RequestIDInterceptor(),
			mw.AuthnBypasswInterceptor(),
			mw.DefaulterInterceptor(),
			mw.ValidatorInterceptor(genericvalidation.NewValidator(c.val)),
		),
	}

	grpcsrv, err := server.NewGRPCServer(
		c.cfg.GRPCOptions,
		serverOptions,
		func(s grpc.ServiceRegistrar) {
			apiv1.RegisterMiniBlogServer(s, handler.NewHandler(c.biz))
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
