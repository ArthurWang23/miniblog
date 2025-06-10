// Copyright 2025 ArthurWang &lt;2826979176@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/arthurwang23/miniblog. The professional
// version of this repository is https://github.com/arthurwang23/miniblog.

package gin

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CORS 跨域资源共享
// 在CORS中HTTP被分为两类：简单请求和复杂请求
// 简单请求：
// 1. 请求方法是GET、HEAD、POST
// 2. 请求头信息不超出以下字段：
//    Accept
//    Accept-Language
//    Content-Language
//    Last-Event-ID
//    Content-Type: 只限于三个值application/x-www-form-urlencoded、multipart/form-data、text/plain

// 复杂请求：
// 凡是不符合简单请求定义的请求

// 浏览器支持CORS功能，当检测到AJAX请求跨域时，会自动添加一些头信息或者添加一次预检请求
// 因此只需服务器实现CORS接口（在HTTP响应头中设置 Access-Control-Allow-Origin）
// 对于简单请求，浏览器直接发出CORS请求，会在请求头信息中添加一个Origin字段，服务器处理头部并在返回头中填充Access-Control-Allow-Origin字段
// 对于复杂请求的CORS跨域处理
// 会在正式通信前添加一次HTTP查询请求（预检请求），使用HTTP方法是OPTIONS 表示请求用于询问目标资源是否允许跨域访问
// 后端收到预检请求后，通过设置跨域相关的HTTP头以完成跨域请求
// miniblog所有请求均为复杂请求，因此只处理复杂请求跨域
// 处理CORS请求中间件
func Cors(c *gin.Context) {
	// 若不是OPTIONS类型的跨域请求，则正常处理该HTTP请求
	// 若为OPTIONS类型的跨域请求，则设置跨域相关的HTTP头，并直接返回响应不再进入后续处理流程
	if c.Request.Method == http.MethodOptions {
		// 必选 设置允许访问的域名 *表示可接受所有域名
		c.Header("Access-Control-Allow-Origin", "*")
		// 必选 逗号分隔的字符串表明支持的所有跨域请求方法
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		// 表示服务器支持的所有头信息字段，不限于浏览器在“预检”中请求的字段，若浏览器请求包含Access-Control-Request-Headers字段，则必选
		c.Header("Access-Control-Allow-Headers", "authorization,origin,content-type,accept")

		c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")

		c.Header("Content-Type", "application/json")

		c.AbortWithStatus(http.StatusOK)

		return
	}
	c.Next()
}

// 用于禁止客户端缓存HTTP请求的返回结果
func NoCache(c *gin.Context) {
	c.Header("Cache-Control", "no-cache,no-store,max-age=0,must-revalidate")
	c.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	c.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	c.Next()
}

// 添加安全相关的HTTP头
func Secure(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("X-Frame-Options", "DENY")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("X-XSS-Protection", "1;mode=block")
	if c.Request.TLS != nil {
		c.Header("Strict-Transport-Security", "max-age=31536000")
	}
	c.Next()
}
