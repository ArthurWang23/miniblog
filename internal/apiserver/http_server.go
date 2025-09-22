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
	httpsrv := server.NewHTTPServer(c.cfg.HTTPOptions, c.cfg.TLSOptions, engin)
	return &ginServer{
		srv: httpsrv,
	}
}

func (c *ServerConfig) InstallRESTAPI(engin *gin.Engine) {
	InstallGenericAPI(engin)

	handler := handler.NewHandler(c.biz, c.val)

	engin.GET("/healthz", handler.Healthz)
	// 这两个接口比较简单，没有API版本
	engin.POST("/login", handler.Login)
	engin.PUT("/refresh-token", mw.AuthnMiddleware(c.retriever), handler.RefreshToken)

	authMiddlewares := []gin.HandlerFunc{mw.AuthnMiddleware(c.retriever), mw.AuthzMiddleware(c.authz)}
	v1 := engin.Group("/v1")
	{
		userv1 := v1.Group("/users")
		{
			userv1.POST("", handler.CreateUser)
			userv1.Use(authMiddlewares...)
			userv1.GET(":userID", handler.GetUser)
			userv1.PUT(":userID", handler.UpdateUser)
			userv1.DELETE(":userID", handler.DeleteUser)
			userv1.PUT(":userID/change-password", handler.ChangePassword)
			userv1.GET("", handler.ListUser)
		}
		postv1 := v1.Group("/posts", authMiddlewares...)
		{
			postv1.POST("", handler.CreatePost)
			postv1.PUT(":postID", handler.UpdatePost)
			postv1.DELETE("", handler.DeletePost)
			postv1.GET(":postID", handler.GetPost)
			postv1.GET("", handler.ListPost)
		}
	}

}

func (c *ServerConfig) InstallRESTAPIWithoutAuth(engin *gin.Engine) {
	InstallGenericAPI(engin)

	handler := handler.NewHandler(c.biz, c.val)

	engin.GET("/healthz", handler.Healthz)
	// 这两个接口比较简单，没有API版本
	engin.POST("/login", handler.Login)
	// 取消 Gin 的鉴权中间件，由 Kratos middleware 负责
	engin.PUT("/refresh-token", handler.RefreshToken)

	v1 := engin.Group("/v1")
	{
		userv1 := v1.Group("/users")
		{
			userv1.POST("", handler.CreateUser)
			// userv1.Use(authMiddlewares...) // 移除 Gin 鉴权
			userv1.GET(":userID", handler.GetUser)
			userv1.PUT(":userID", handler.UpdateUser)
			userv1.DELETE(":userID", handler.DeleteUser)
			userv1.PUT(":userID/change-password", handler.ChangePassword)
			userv1.GET("", handler.ListUser)
		}
		postv1 := v1.Group("/posts")
		{
			postv1.POST("", handler.CreatePost)
			postv1.PUT(":postID", handler.UpdatePost)
			postv1.DELETE("", handler.DeletePost)
			postv1.GET(":postID", handler.GetPost)
			postv1.GET("", handler.ListPost)
		}
	}
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
