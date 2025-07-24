package store

import "github.com/ArthurWang23/miniblog/internal/pkg/log"

type Logger struct{}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) Error(err error, msg string, kvs ...any) {
	log.Errorw(msg, append(kvs, "err", err)...)
}
