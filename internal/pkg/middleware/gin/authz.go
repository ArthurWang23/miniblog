package gin

import (
	"github.com/ArthurWang23/miniblog/internal/pkg/contextx"
	"github.com/ArthurWang23/miniblog/internal/pkg/errno"
	"github.com/ArthurWang23/miniblog/internal/pkg/log"
	"github.com/ArthurWang23/miniblog/pkg/core"
	"github.com/gin-gonic/gin"
)

type Authorizer interface {
	Authorize(subject, object, action string) (bool, error)
}

func AuthzMiddleware(authorizer Authorizer) gin.HandlerFunc {
	return func(c *gin.Context) {
		subject := contextx.UserID(c.Request.Context())
		object := c.Request.URL.Path
		action := c.Request.Method

		log.Debugw("Build authorize context", "subject", subject, "object", object, "action", action)
		if allowed, err := authorizer.Authorize(subject, object, action); err != nil || !allowed {
			core.WriteResponse(c, nil, errno.ErrPermissionDenied.WithMessage(
				"access denied : subject=%s,object=%s,action=%s,reason=%v",
				subject,
				object,
				action,
				err,
			))
			c.Abort()
			return
		}
		c.Next()
	}
}
