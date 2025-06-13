// Copyright 2025 ArthurWang &lt;2826979176@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/arthurwang23/miniblog. The professional
// version of this repository is https://github.com/arthurwang23/miniblog.

package contextx

import "context"

type (
	// userIDKey 定义用户ID的上下文键
	userIDKey struct{}

	usernameKey struct{}

	accessTokenKey struct{}
	// requestIDKey 定义请求ID的上下文键
	requestIDKey struct{}
)

// 将userID存放到上下文中
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey{}, userID)
}

// 从上下文提取id
func UserID(ctx context.Context) string {
	userID, _ := ctx.Value(userIDKey{}).(string)
	return userID
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, requestID)
}

func RequestID(ctx context.Context) string {
	requestID, _ := ctx.Value(requestIDKey{}).(string)
	return requestID
}

func WithUsername(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, usernameKey{}, username)
}

func Username(ctx context.Context) string {
	username, _ := ctx.Value(usernameKey{}).(string)
	return username
}

func WithAccessToken(ctx context.Context, accessToken string) context.Context {
	return context.WithValue(ctx, accessTokenKey{}, accessToken)
}

func AccessToken(ctx context.Context) string {
	accessToken, _ := ctx.Value(accessTokenKey{}).(string)
	return accessToken
}
