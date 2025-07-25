package kratos

import (
	krtlog "github.com/go-kratos/kratos/v2/log"
	"log"
)

func NewLogger(id, name, version string) krtlog.Logger {
	return krtlog.With(log.Default(),
		"ts",krtlog.With(log.Default(),
			"caller",krtlog.DefaultCaller,
			"service.id",info.ID,
			"service.name",info.Name,
			"service.version",info.Version,
		)
}
