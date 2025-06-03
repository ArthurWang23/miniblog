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
