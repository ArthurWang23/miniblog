// Copyright 2025 ArthurWang &lt;2826979176@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/arthurwang23/miniblog. The professional
// version of this repository is https://github.com/arthurwang23/miniblog.

package server

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"

	"github.com/ArthurWang23/miniblog/internal/pkg/log"
	genericoptions "github.com/ArthurWang23/miniblog/pkg/options"
)

type HTTPServer struct {
	srv *http.Server
}

func NewHTTPServer(httpOptions *genericoptions.HTTPOptions, tlsOptions *genericoptions.TLSOptions, handler http.Handler) *HTTPServer {
	var tlsConfig *tls.Config
	if tlsOptions != nil && tlsOptions.UseTLS {
		tlsConfig = tlsOptions.MustTLSConfig()
	}

	return &HTTPServer{
		srv: &http.Server{
			Addr:      httpOptions.Addr,
			TLSConfig: tlsConfig,
			Handler:   handler,
		},
	}
}

func (s *HTTPServer) RunOrDie() {
	log.Infow("Start to listening the incoming requests", "protocol", protocolName(s.srv), "addr", s.srv.Addr)
	// 默认启用http
	serveFn := func() error { return s.srv.ListenAndServe() }
	if s.srv.TLSConfig != nil {
		serveFn = func() error { return s.srv.ListenAndServeTLS("", "") }
	}
	if err := serveFn(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalw("Failed to server HTTP(s) server", "err", err)
	}
}

func (s *HTTPServer) GracefulStop(ctx context.Context) {
	log.Infow("Gracefully stop HTTP(s) server")
	// Shutdown 方法首先关闭所有已开启的监听器，然后关闭所有空闲连接
	// 最后等待所有活跃连接进入空闲状态后终止服务
	// 若传入ctx在服务完成终止之前超时，则Shutdown会返回与context相关的错误
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Errorw("HTTP(s) server forced to shutdown", "err", err)
	}
}
