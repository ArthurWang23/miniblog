package http

import (
	"github.com/ArthurWang23/miniblog/pkg/core"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreatePost(c *gin.Context) {
	core.HandleJSONRequest(c, h.biz.PostV1().Create)
}

func (h *Handler) UpdatePost(c *gin.Context) {
	core.HandleJSONRequest(c, h.biz.PostV1().Update)
}

func (h *Handler) DeletePost(c *gin.Context) {
	core.HandleJSONRequest(c, h.biz.PostV1().Delete)
}

func (h *Handler) GetPost(c *gin.Context) {
	core.HandleUriRequest(c, h.biz.PostV1().Get)
}

func (h *Handler) ListPost(c *gin.Context) {
	core.HandleQueryRequest(c, h.biz.PostV1().List)
}
