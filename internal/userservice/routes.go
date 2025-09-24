package userservice

import (
	handler "github.com/ArthurWang23/miniblog/internal/userservice/handler/http"
	"github.com/gin-gonic/gin"
)

func InstallRoutes(engine *gin.Engine, h *handler.Handler) {
	h.Register(engine)
}