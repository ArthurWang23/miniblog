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


// 在开发go共享包时，要遵循包功能完整，稳定，独立，可定制化原则
// 通常使用函数选项模式
// 不要使用github.com/ArthurWang23/miniblog/internal/pkg/log
// 这种项目定制的日志包，因为不同项目使用的日志包是不一样的
// 可以通过WithLogger函数选项来设置调用发使用的Logger
// 共享包还需要避免使用init，panic 这种调用方很难感知的代码实现


// server包是工呢个独立完整的共享包，不会实现业务相关代码
// 在internal/apiserver要基于server包进一步结构化代码
// 详见internal/apiserver/server.go

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
