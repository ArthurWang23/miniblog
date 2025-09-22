package log

import (
	"fmt"
	krtlog "github.com/go-kratos/kratos/v2/log"
)

type KratosLogger interface {
	Log(level krtlog.Level, keyvals ...any) error
}

func (l *zapLogger) Log(level krtlog.Level, keyvals ...any) error {
	keylen := len(keyvals)
	if keylen == 0 || keylen%2 != 0 {
		l.z.Warn(fmt.Sprintf("Keyvalues must appear in pairs: %v", keyvals))
		return nil
	}

	switch level {
	case krtlog.LevelDebug:
		l.z.Sugar().Debugw("", keyvals...)
	case krtlog.LevelInfo:
		l.z.Sugar().Infow("", keyvals...)
	case krtlog.LevelWarn:
		l.z.Sugar().Warnw("", keyvals...)
	case krtlog.LevelError:
		l.z.Sugar().Errorw("", keyvals...)
	case krtlog.LevelFatal:
		l.z.Sugar().Fatalw("", keyvals...)
	}
	return nil
}

// 提供一个便捷方法，直接以 krtlog.Logger 形式返回
func Kratos() krtlog.Logger {
	return std
}
