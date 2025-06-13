package http

import (
	"github.com/ArthurWang23/miniblog/pkg/core"
	"github.com/gin-gonic/gin"
)

// http 使用了gin web框架，所以在解析请求和返回请求时需要使用gin框架提供的各类方法
// core.HandleJSONRequest 函数封装了请求参数解析、参数校验、请求返回逻辑
func (h *Handler) Login(c *gin.Context) {
	core.HandleJSONRequest(c, h.biz.UserV1().Login)
}

func (h *Handler) RefreshToken(c *gin.Context) {
	core.HandleJSONRequest(c, h.biz.UserV1().RefreshToken)
}

func (h *Handler) ChangePassword(c *gin.Context) {
	core.HandleJSONRequest(c, h.biz.UserV1().ChangePassword)
}

func (h *Handler) CreateUser(c *gin.Context) {
	core.HandleJSONRequest(c, h.biz.UserV1().Create)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	core.HandleJSONRequest(c, h.biz.UserV1().Update)
}

func (h *Handler) DeleteUser(c *gin.Context) {
	core.HandleUriRequest(c, h.biz.UserV1().Delete)
}

func (h *Handler) GetUser(c *gin.Context) {
	core.HandleUriRequest(c, h.biz.UserV1().Get)
}

func (h *Handler) ListUser(c *gin.Context) {
	core.HandleQueryRequest(c, h.biz.UserV1().List)
}
