package kratos

import (
	"os"

	krtlog "github.com/go-kratos/kratos/v2/log"
)

// NewLogger 返回一个带有标准字段的 Kratos Logger。
func NewLogger(id, name, version string) krtlog.Logger {
	return krtlog.With(
		krtlog.NewStdLogger(os.Stdout),
		"ts", krtlog.DefaultTimestamp,
		"caller", krtlog.DefaultCaller,
		"service.id", id,
		"service.name", name,
		"service.version", version,
	)
}
