package grpc

import (
	"context"

	"google.golang.org/grpc"
)

type RequestValidator interface {
	Validate(ctx context.Context, rq any) error
}

func ValidatorInterceptor(validator RequestValidator) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, rq any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if err := validator.Validate(ctx, rq); err != nil {
			return nil, err
		}
		return handler(ctx, rq)
	}
}
