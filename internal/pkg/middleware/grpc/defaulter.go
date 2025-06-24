package grpc

import (
	"context"

	"google.golang.org/grpc"
)

func DefaulterInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, rq any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		// 有default方法则调用
		if defaulter, ok := rq.(interface{ Default() }); ok {
			defaulter.Default()
		}
		return handler(ctx, rq)
	}
}
