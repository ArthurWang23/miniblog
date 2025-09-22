package krt

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware"
)

type RequestValidator interface {
	Validate(ctx context.Context, rq any) error
}

func Validator(v RequestValidator) middleware.Middleware {
	return func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, rq interface{}) (interface{}, error) {
			if err := v.Validate(ctx, rq); err != nil {
				return nil, err
			}
			return next(ctx, rq)
		}
	}
}