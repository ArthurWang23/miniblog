// Copyright 2025 ArthurWang &lt;2826979176@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/arthurwang23/miniblog. The professional
// version of this repository is https://github.com/arthurwang23/miniblog.

// 基于cobra构建go框架
// 为了代码清晰，main中只有核心启动的代码
package main

import (
	"os"

	"github.com/ArthurWang23/miniblog/cmd/mb-apiserver/app"

	// 自动设置GOMAXPROCS
	_ "go.uber.org/automaxprocs"
)

// Go 程序的默认入口函数。阅读项目代码的入口函数.
func main() {
	command := app.NewMiniBlogCommand()

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
