package server

import (
	"context"
	"net/http"
)

// 重构server.go代码
// grpc服务器、http反向代理服务器、gin http服务器
// 均需要以下方法：
// NewXXX 创建一个XXX服务器类型
// Run 启动XXX服务器
// GracefulStop 优雅关停XXX服务器的方法

// 抽象Server接口类型

type Server interface {
	RunOrDie()

	GracefulStop(ctx context.Context)
}

func protocolName(server *http.Server) string {
	if server.TLSConfig != nil {
		return "https"
	}
	return "http"
}
