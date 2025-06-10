// Copyright 2025 ArthurWang &lt;2826979176@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/arthurwang23/miniblog. The professional
// version of this repository is https://github.com/arthurwang23/miniblog.

package known

// 将一些共享常量统一保存在常量包 如known constant这类包中，以便集中管理和引用

const (
	// XRequestID 用来定义上下文中的键，代表请求ID
	XRequestID = "x-request-id"

	// XUserID 用来定义上下文的键，代表请求用户ID UserID整个用户生命周期唯一
	XUserID = "x-user-id"
)
