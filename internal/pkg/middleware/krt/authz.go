package krt

import (
	"context"

	"github.com/ArthurWang23/miniblog/internal/pkg/contextx"
	"github.com/ArthurWang23/miniblog/internal/pkg/errno"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
)

type Authorizer interface {
	Authorize(subject, object, action string) (bool, error)
}

func Authz(authorizer Authorizer, whitelist map[string]struct{}) middleware.Middleware {
	return func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if info, ok := transport.FromServerContext(ctx); ok {
				// 统一计算 object 与 action
				object := info.Operation()
				action := "CALL"

				// HTTP 场景：使用 Path + Method
				if ht, ok := info.(*kratoshttp.Transport); ok && ht.Request() != nil {
					r := ht.Request()
					object = r.URL.Path
					action = r.Method
				}

				// 白名单（gRPC: FullMethod；HTTP: Path）
				if _, skip := whitelist[object]; skip {
					return next(ctx, req)
				}
				subject := contextx.UserID(ctx)

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