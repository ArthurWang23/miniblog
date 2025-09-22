package krt

import (
	"context"

	"github.com/ArthurWang23/miniblog/internal/pkg/contextx"
	"github.com/ArthurWang23/miniblog/internal/pkg/errno"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

type Authorizer interface {
	Authorize(subject, object, action string) (bool, error)
}

func Authz(authorizer Authorizer, whitelist map[string]struct{}) middleware.Middleware {
	return func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if info, ok := transport.FromServerContext(ctx); ok {
				// 白名单方法跳过
				if _, skip := whitelist[info.Operation()]; skip {
					return next(ctx, req)
				}
				subject := contextx.UserID(ctx)
				object := info.Operation()
				action := "CALL"

				log.Debugf("Build authorize context", "subject", subject, "object", object, "action", action)

				if allowed, err := authorizer.Authorize(subject, object, action); err != nil || !allowed {
					return nil, errno.ErrPermissionDenied.WithMessage(
						"access denied : subject=%s,object=%s,action=%s,reason=%v",
						subject, object, action, err,
					)
				}
			}
			return next(ctx, req)
		}
	}
}