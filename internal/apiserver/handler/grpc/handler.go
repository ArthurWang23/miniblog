// Copyright 2025 ArthurWang &lt;2826979176@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/arthurwang23/miniblog. The professional
// version of this repository is https://github.com/arthurwang23/miniblog.

package grpc

import (
	"github.com/ArthurWang23/miniblog/internal/apiserver/biz"
	apiv1 "github.com/ArthurWang23/miniblog/pkg/api/apiserver/v1"
)

type Handler struct {
	// 提供默认实现，确保未实现的grpc方法返回未实现错误
	// 确保向后兼容，当接口新增方法时，服务端实现不需要立即定义新方法
	// 简化服务实现过程，开发者只需要实现自己需要的方法，不必为每个方法提供一个默认的未实现错误
	// 提高代码可维护性
	apiv1.UnimplementedMiniBlogServer

	biz biz.IBiz
}

func NewHandler(biz biz.IBiz) *Handler {
	return &Handler{
		biz: biz,
	}
}
