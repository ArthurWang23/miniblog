// Copyright 2025 ArthurWang &lt;2826979176@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/arthurwang23/miniblog. The professional
// version of this repository is https://github.com/arthurwang23/miniblog.

package grpc

import (
	"context"

	"github.com/ArthurWang23/miniblog/internal/pkg/contextx"
	"github.com/ArthurWang23/miniblog/internal/pkg/known"
	"github.com/ArthurWang23/miniblog/pkg/errorsx"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// grpc拦截器
// 为每次请求注入一个唯一的 RequestID，并在每条日志中输出该 RequestID
// 用户只需提供一个唯一的 RequestID，开发人员即可快速定位与该请求相关的所有日志记录
// 具体实现在请求中注入 RequestID  在日志中打印 RequestID
// 因为每个请求都需要requestid 因此使用grpc拦截器实现

// 针对服务端一元调用（Unary RPC）的服务端一元拦截器RequestIDInterceptor
// 返回的函数先尝试从grpc请求的元数据中获取键为x-request-id的请求id，若没有获取则调用uuid创建id
// 并将id保存到grpc请求和返回的元数据中

// 元数据
// grpc中，元数据时一种轻量级、灵活的机制，用于通过上下文传递额外的信息
// 如认证信息，追踪id等，类似于http中的header
// 元数据本质上是一组键值对，可以在rpc调用的请求和响应双方进行交换
// 这些键值通常被用作元信息，而不是直接与业务数据相关

func RequestIDInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		var requestID string
		md, _ := metadata.FromIncomingContext(ctx)
		if requestIDs := md[known.XRequestID]; len(requestIDs) > 0 {
			requestID = requestIDs[0]
		}
		// 如果请求id为空，则生成一个唯一的请求id
		if requestID == "" {
			requestID = uuid.New().String()
			md.Append(known.XRequestID, requestID)
		}

		// 将元数据设为新的incoming context
		ctx = metadata.NewIncomingContext(ctx, md)
		// 将请求ID设置到响应的Header Metadata中
		// SetHeader在grpc方法响应中添加元数据
		// 仅设置数据，不会立即发送给客户端
		// header metadata会在rpc响应返回时一并发送
		grpc.SetHeader(ctx, md)

		// 将请求 ID 添加到自定义的上下文中
		// grpc请求的返回元数据过程中,已经包含了请求id
		// 这里再次将请求id添加到grpc返回错误中，原因如下：
		// 错误中包含请求id更易定位问题
		// grpc元数据很少被处理

		ctx = contextx.WithRequestID(ctx, requestID)

		res, err := handler(ctx, req)
		if err != nil {
			return res, errorsx.FromError(err).WithRequestID(requestID)
		}
		return res, nil
	}
}
