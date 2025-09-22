package krt

import (
	"context"

	"github.com/ArthurWang23/miniblog/internal/pkg/contextx"
	"github.com/ArthurWang23/miniblog/internal/pkg/known"
	"github.com/ArthurWang23/miniblog/pkg/errorsx"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/google/uuid"
)

func RequestID() middleware.Middleware {
	return func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			requestID := ""
			if info, ok := transport.FromServerContext(ctx); ok {
				// 优先从请求头获取
				if v := info.RequestHeader().Get(known.XRequestID); v != "" {
					requestID = v
				} else {
					// 不存在则生成一个并写入响应头
					requestID = uuid.New().String()
					info.ReplyHeader().Set(known.XRequestID, requestID)
				}
			} else {
				// 兜底生成
				requestID = uuid.New().String()
			}

			// 注入到上下文，便于后续日志与错误关联
			ctx = contextx.WithRequestID(ctx, requestID)

			res, err := next(ctx, req)
			if err != nil {
				return res, errorsx.FromError(err).WithRequestID(requestID)
			}
			return res, nil
		}
	}
}