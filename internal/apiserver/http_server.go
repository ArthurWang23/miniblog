package apiserver

import (
	"context"
	"net/http"

	handler "github.com/ArthurWang23/miniblog/internal/apiserver/handler/http"
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
