// Copyright 2025 ArthurWang &lt;2826979176@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/arthurwang23/miniblog. The professional
// version of this repository is https://github.com/arthurwang23/miniblog.

package server

import (
	"context"
	"net"

	"github.com/ArthurWang23/miniblog/internal/pkg/log"
	genericoptions "github.com/ArthurWang23/miniblog/pkg/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	srv *grpc.Server
	lis net.Listener
}

func NewGRPCServer(
	grpcOptions *genericoptions.GRPCOptions,
	serverOptions []grpc.ServerOption,
	registerServer func(grpc.ServiceRegistrar),
) (*GRPCServer, error) {
	lis, err := net.Listen("tcp", grpcOptions.Addr)
	if err != nil {
		log.Errorw("Failed to listen", "err", err)
		return nil, err
	}
	grpcsrv := grpc.NewServer(serverOptions...)
	registerServer(grpcsrv)
	registerHealthServer(grpcsrv)
	reflection.Register(grpcsrv)
	return &GRPCServer{
		srv: grpcsrv,
		lis: lis,
	}, nil
}

func (s *GRPCServer) RunOrDie() {
	log.Infow("Start to listening the incoming requests", "protocol", "grpc", "addr", s.lis.Addr().String())
	if err := s.srv.Serve(s.lis); err != nil {
		log.Fatalw("Failed to serve grpc server", "err", err)
	}
}

func (s *GRPCServer) GracefulStop(ctx context.Context) {
	log.Infow("Gracefully stop grpc server")
	s.srv.GracefulStop()
}

// 注册健康检查服务
func registerHealthServer(grpcsrv *grpc.Server) {
	healthServer := health.NewServer()
	healthServer.SetServingStatus("MiniBlog", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(grpcsrv, healthServer)
}
