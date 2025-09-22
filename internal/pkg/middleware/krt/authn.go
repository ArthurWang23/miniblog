package krt

import (
	"context"

	"github.com/ArthurWang23/miniblog/internal/apiserver/model"
	"github.com/ArthurWang23/miniblog/internal/pkg/contextx"
	"github.com/ArthurWang23/miniblog/internal/pkg/errno"
	"github.com/ArthurWang23/miniblog/internal/pkg/known"
	"github.com/ArthurWang23/miniblog/internal/pkg/log"
	"github.com/ArthurWang23/miniblog/pkg/token"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
)

type UserRetriever interface {
	GetUser(ctx context.Context, userID string) (*model.UserM, error)
}

func Authn(retriever UserRetriever, whitelist map[string]struct{}) middleware.Middleware {
	return func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			// 白名单方法/路径直接跳过认证
			if info, ok := transport.FromServerContext(ctx); ok {
				target := info.Operation()
				if ht, ok := info.(*kratoshttp.Transport); ok && ht.Request() != nil {
					target = ht.Request().URL.Path
				}
				if _, skip := whitelist[target]; skip {
					return next(ctx, req)
				}
			}

			userID, err := token.ParseRequest(ctx)
			if err != nil {
				log.Errorw("Failed to parse request", "err", err)
				return nil, errno.ErrTokenInvalid.WithMessage(err.Error())
			}
			log.Debugw("Token parsing successful", "userID", userID)

			user, err := retriever.GetUser(ctx, userID)
			if err != nil {
				return nil, errno.ErrUnauthenticated.WithMessage(err.Error())
			}

			// 注入用户信息到上下文
			ctx = context.WithValue(ctx, known.XUsername, user.Username)
			ctx = context.WithValue(ctx, known.XUserID, userID)
			ctx = contextx.WithUserID(ctx, user.UserID)
			ctx = contextx.WithUsername(ctx, user.Username)

			return next(ctx, req)
		}
	}
}