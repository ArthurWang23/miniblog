package grpc

import (
	"context"

	"github.com/ArthurWang23/miniblog/internal/pkg/contextx"
	"github.com/ArthurWang23/miniblog/internal/pkg/known"
	"github.com/ArthurWang23/miniblog/internal/pkg/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// 在开发了 Store 层、Biz 层、Handler 层代码之后，便可以对整个项目代码进行初步的测试
// 以尽快验证代码的设计是否合理、核心功能是否可用。
// 因为租户 UserID 数据是从请求的 Token 中获取的，这时候整个项目还未实现认证功能
// 为了能够测试项目代码，可以开发一个 bypass 认证中间件
// bypass 认证中间件会从请求头中获取用户的 UserID 数据，并放通所有请求。
// 通过 bypass 中间件，既能够获取到租户数据，又能够让请求认证通过。

// grpc拦截器，模拟所有请求都通过认证
func AuthnBypasswInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		userID := "user-000001" // 默认用户id
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if values := md.Get(known.XUserID); len(values) > 0 {
				userID = values[0]
			}
		}
		log.Debugw("Simulated authentication successful", "userID", userID)

		ctx = context.WithValue(ctx, known.XUserID, userID)
		// 为log和contextx提供用户上下文支持
		ctx = contextx.WithUserID(ctx, userID)

		return handler(ctx, req)
	}
}
