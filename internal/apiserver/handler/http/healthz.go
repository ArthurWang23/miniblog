package http

import (
	"time"

	apiv1 "github.com/ArthurWang23/miniblog/pkg/api/apiserver/v1"
	"github.com/gin-gonic/gin"
)

func (h *Handler) Healthz(c *gin.Context) {
	c.JSON(200, &apiv1.HealthzResponse{
		Status:    apiv1.ServiceStatue_Healthy,
		Timestamp: time.Now().Format(time.DateTime),
	})
}
