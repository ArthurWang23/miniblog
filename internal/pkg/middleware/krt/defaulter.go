package krt

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware"
)

func Defaulter() middleware.Middleware {
	return func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, rq interface{}) (interface{}, error) {
			if defaulter, ok := rq.(interface{ Default() }); ok {
				defaulter.Default()
			}
			return next(ctx, rq)
		}
	}
}