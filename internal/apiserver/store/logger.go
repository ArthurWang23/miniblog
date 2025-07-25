package store

import (
	"context"
	"github.com/ArthurWang23/miniblog/internal/pkg/log"
)

type Logger struct{}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) Error(ctx context.Context, err error, msg string, kvs ...any) {
	log.W(ctx).Errorw(msg, append(kvs, "err", err)...)
}
