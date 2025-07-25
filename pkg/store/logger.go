package store

import "context"

type Logger interface {
	Error(ctx context.Context, err error, message string, kvs ...any)
}
