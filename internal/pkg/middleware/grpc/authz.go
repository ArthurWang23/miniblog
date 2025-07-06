package grpc

import (
	"context"

	"github.com/ArthurWang23/miniblog/internal/pkg/contextx"
	"github.com/ArthurWang23/miniblog/internal/pkg/errno"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc"
)

type Authorizer interface {
	Authorize(subject, object, action string) (bool, error)
}

func AuthzInterceptor(authorizer Authorizer) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		subject := contextx.UserID(ctx)
		object := info.FullMethod
		action := "CALL"

		log.Debugf("Build authorize context", "subject", subject, "object", object, "action", action)

		if allowed, err := authorizer.Authorize(subject, object, action); err != nil || !allowed {
			return nil, errno.ErrPermissionDenied.WithMessage(
				"access denied : subject=%s,object=%s,action=%s,reason=%v",
				subject,
				object,
				action,
				err,
			)
		}
		return handler(ctx, req)
	}
}
