// Copyright 2025 ArthurWang &lt;2826979176@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/arthurwang23/miniblog. The professional
// version of this repository is https://github.com/arthurwang23/miniblog.

package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/ArthurWang23/miniblog/internal/pkg/log"
	genericoptions "github.com/ArthurWang23/miniblog/pkg/options"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

type GRPCGatewayServer struct {
	srv *http.Server
}

func NewGRPCGatewayServer(
	httpOptions *genericoptions.HTTPOptions,
	grpcOptions *genericoptions.GRPCOptions,
	registerHandler func(mux *runtime.ServeMux, conn *grpc.ClientConn) error,
) (*GRPCGatewayServer, error) {
	dialOptions := []grpc.DialOption{
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff:           backoff.DefaultConfig,
			MinConnectTimeout: 10 * time.Second,
		}),
	}
	dialOptions = append(dialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(grpcOptions.Addr, dialOptions...)
	if err != nil {
		log.Errorw("Failed to dial context", "err", err)
		return nil, err
	}
	gwmux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseEnumNumbers: true,
		},
	}))
	if err := registerHandler(gwmux, conn); err != nil {
		log.Errorw("Failed to register handler", "err", err)
		return nil, err
	}

	return &GRPCGatewayServer{
		srv: &http.Server{
			Addr:    httpOptions.Addr,
			Handler: gwmux,
		},
	}, nil
}

// 启动GRPC网关服务器并在出错时记录致命错误
func (s *GRPCGatewayServer) RunOrDie() {
	log.Infow("Start to listening the incoming requests", "protocol", protocolName(s.srv), "addr", s.srv.Addr)
	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalw("Failed to server HTTP(s) server", "err", err)
	}
}

func (s *GRPCGatewayServer) GracefulStop(ctx context.Context) {
	log.Infow("Gracefully stop HTTP(s) server")
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Errorw("HTTP(s) server forced to shutdown", "err", err)
	}
}
