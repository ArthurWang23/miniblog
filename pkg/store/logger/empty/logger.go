package empty

import "context"

// a no-op logger that implements the Logger interface

type emptyLogger struct{}

func NewLogger() *emptyLogger {
	return &emptyLogger{}
}

func (l *emptyLogger) Error(cfx context.Context, err error, msg string, kvs ...any) {
	// no operation
}
