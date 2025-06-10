// Copyright 2025 ArthurWang &lt;2826979176@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/arthurwang23/miniblog. The professional
// version of this repository is https://github.com/arthurwang23/miniblog.

package apiserver

import (
	"context"
	"net/http"

	handler "github.com/ArthurWang23/miniblog/internal/apiserver/handler/http"
	mw "github.com/ArthurWang23/miniblog/internal/pkg/middleware/gin"
	"github.com/ArthurWang23/miniblog/internal/pkg/server"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

type ginServer struct {
	srv server.Server
}

var _ server.Server = (*ginServer)(nil)

func (c *ServerConfig) NewGinServer() server.Server {
	engin := gin.New()
	// 先注册中间件，再注册路由
	engin.Use(gin.Recovery(), mw.NoCache, mw.Cors, mw.Secure, mw.RequestIDMiddleware())
	// 注册rest api 路由
	c.InstallRESTAPI(engin)
	httpsrv := server.NewHTTPServer(c.cfg.HTTPOptions, engin)
	return &ginServer{
		srv: httpsrv,
	}
}

func (c *ServerConfig) InstallRESTAPI(engin *gin.Engine) {
	InstallGenericAPI(engin)

	handler := handler.NewHandler()

	engin.GET("/healthz", handler.Healthz)
}

// 安装业务无关的路由 ，如pprof、404
func InstallGenericAPI(engin *gin.Engine) {
	// pprof用来提供性能调试和优化的api接口
	pprof.Register(engin)

	engin.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, "Page not found.")
	})
}

func (s *ginServer) RunOrDie() {
	s.srv.RunOrDie()
}

func (s *ginServer) GracefulStop(ctx context.Context) {
	s.srv.GracefulStop(ctx)
}
