package onex

import "github.com/ArthurWang23/miniblog/pkg/log"

type onexLogger struct{}

func NewLogger() *onexLogger {
	return &onexLogger{}
}

func (l *onexLogger) Error(err error, msg string, kvs ...any) {
	log.Errorw(err, msg, kvs)
}
