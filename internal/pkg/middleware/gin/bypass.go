package gin

import (
	"github.com/ArthurWang23/miniblog/internal/pkg/contextx"
	"github.com/ArthurWang23/miniblog/internal/pkg/known"
	"github.com/ArthurWang23/miniblog/internal/pkg/log"
	"github.com/gin-gonic/gin"
)

func AuthnBypasswMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := "user-000001"
		if val := c.GetHeader(known.XUserID); val != "" {
			userID = val
		}
		log.Debugw("Simulated anthentication successful", "userID", userID)

		ctx := contextx.WithUserID(c.Request.Context(), userID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
