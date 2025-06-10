// Copyright 2025 ArthurWang &lt;2826979176@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/arthurwang23/miniblog. The professional
// version of this repository is https://github.com/arthurwang23/miniblog.

package main

import (
	"fmt"

	"github.com/ArthurWang23/miniblog/internal/pkg/errno"
	"github.com/ArthurWang23/miniblog/pkg/errorsx"
)

func main() {
	errx := errorsx.New(500, "InternalError.DBConnection", "Something went wrong: %s", "DB connection failed")
	// 调用errx的Error方法
	fmt.Println(errx)

	errx.WithMetadata(map[string]string{
		"user_id":    "12345",
		"request_id": "abc-def",
	})

	errx.KV("trace_id", "xyz-789")

	errx.WithMessage("Updated message:%s", "retry failed")

	fmt.Println(errx)

	someerr := doSomething()

	fmt.Println(someerr)

	fmt.Println(errno.ErrUsernameInvalid.Is(someerr))

	fmt.Println(errno.ErrPasswordInvalid.Is(someerr))

}

func doSomething() error {
	return errno.ErrUsernameInvalid.WithMessage("Username is too short")
}
