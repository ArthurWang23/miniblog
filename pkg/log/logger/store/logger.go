package store

import (
	"context"
	"github.com/ArthurWang23/miniblog/pkg/log"
)

type Logger struct{}

// NewLogger creates and returns a new instance of Logger.
func NewLogger() *Logger {
	return &Logger{}
}

// Error logs an error message with the provided context using the log package.
func (l *Logger) Error(ctx context.Context, err error, msg string, kvs ...any) {
	log.W(ctx).Errorw(err, msg, kvs...)
}
