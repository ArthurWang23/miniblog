//go:build wireinject
// +build wireinject

// go:build wireinject
package apiserver

import (
	"github.com/ArthurWang23/miniblog/internal/apiserver/biz"
	"github.com/ArthurWang23/miniblog/internal/apiserver/pkg/validation"
	"github.com/ArthurWang23/miniblog/internal/apiserver/store"
	ginmw "github.com/ArthurWang23/miniblog/internal/pkg/middleware/gin"
	"github.com/ArthurWang23/miniblog/internal/pkg/server"
	"github.com/ArthurWang23/miniblog/pkg/auth"
	"github.com/google/wire"
)

// 通过wire实现依赖注入
func InitializeWebServer(*Config) (server.Server, error) {
	wire.Build(
		wire.NewSet(NewWebServer, wire.FieldsOf(new(*Config), "ServerMode")),
		wire.Struct(new(ServerConfig), "*"), // 表示注入全部字段
		wire.NewSet(store.ProviderSet, biz.ProviderSet),
		ProviderDB,
		validation.ProviderSet,
		wire.NewSet(
			wire.Struct(new(UserRetriever), "*"),
			wire.Bind(new(ginmw.UserRetriever), new(*UserRetriever)),
		),
		auth.ProviderSet,
	)
	return nil, nil
}
