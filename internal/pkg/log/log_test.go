package log_test

import (
	"testing"
	"time"

	"github.com/ArthurWang23/miniblog/internal/pkg/log"
)

func TestLogger(t *testing.T) {
	opts := &log.Options{
		Level:             "debug",
		Format:            "json",
		DisableCaller:     false,
		DisableStacktrace: false,
		OutputPaths:       []string{"stdout"},
	}

	log.Init(opts)

	log.Debugw("This is a debug message", "key1", "value1", "key2", "value2")
	log.Infow("This is an info message", "key", "value")
	log.Warnw("This is a warning message", "timestamp", time.Now())
	log.Errorw("This is an error message", "error", "something went wrong")
	// 确保日志缓冲区被刷新
	log.Sync()
}
