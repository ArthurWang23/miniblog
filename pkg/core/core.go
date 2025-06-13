package core

import (
	"context"
	"net/http"

	"github.com/ArthurWang23/miniblog/pkg/errorsx"
	"github.com/gin-gonic/gin"
)

// 验证函数的类型，用于对绑定的数据结构进行验证
type Validator[T any] func(context.Context, *T) error

// 定义绑定函数的类型，用于绑定请求数据到相应结构体
type Binder func(any) error

// 处理函数的类型，用于处理已经绑定和验证的数据
type Handler[T any, R any] func(ctx context.Context, req *T) (R, error)

// 用于API请求中发生错误时返回统一的格式化错误信息
type ErrorResponse struct {
	Reason   string            `json:"reason,omitempty"`
	Message  string            `json:"message,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// 处理JSON请求的快捷函数
// 封装HandleRequest的语法糖
func HandleJSONRequest[T any, R any](c *gin.Context, handler Handler[T, R], validators ...Validator[T]) {
	HandleRequest(c, c.ShouldBindJSON, handler, validators...)
}

// 处理Query请求的快捷函数
func HandleQueryRequest[T any, R any](c *gin.Context, handler Handler[T, R], validators ...Validator[T]) {
	HandleRequest(c, c.ShouldBindQuery, handler, validators...)
}

// 处理URI请求的快捷函数
// 在调用 *gin.Context 的 ShouldBindUri 方法时
// Gin 框架会将请求中的路径参数绑定到 Go 结构体中的对应字段上
// 这些字段跟路径参数的映射关系，是通过 Go 结构体字段的 uri 标签来映射的
// protoc 编译器生成的 Go 结构体字段中的标签是不带 uri 标签的
// Go 项目开发中，可以使用 protoc-go-inject-tag 工具来给 Protobuf 消息定义中的指定字段添加 Go 结构体标签

func HandleUriRequest[T any, R any](c *gin.Context, handler Handler[T, R], validators ...Validator[T]) {
	HandleRequest(c, c.ShouldBindUri, handler, validators...)
}

// 通用的请求处理函数
// 负责绑定请求数据，执行验证，并调用实际的业务处理逻辑函数
func HandleRequest[T any, R any](c *gin.Context, binder Binder, handler Handler[T, R], validators ...Validator[T]) {
	var request T
	// 绑定和验证请求数据
	// 先调用ReadRequest从请求中解析出参数，所有的请求信息都保存在*gin.Context的变量c中
	// 解析完参数会调用handler方法
	if err := ReadRequest(c, &request, binder, validators...); err != nil {
		WriteResponse(c, nil, err)
		return
	}
	response, err := handler(c.Request.Context(), &request)

	WriteResponse(c, response, err)
}

func ShouldBindJSON[T any](c *gin.Context, rq *T, validators ...Validator[T]) error {
	return ReadRequest(c, rq, c.ShouldBindJSON, validators...)
}

func ShouldBindQuery[T any](c *gin.Context, rq *T, validators ...Validator[T]) error {
	return ReadRequest(c, rq, c.ShouldBindQuery, validators...)
}

func ShouldBindUri[T any](c *gin.Context, rq *T, validators ...Validator[T]) error {
	return ReadRequest(c, rq, c.ShouldBindUri, validators...)
}

// 用于绑定和验证请求数据的通用工具函数
// 调用绑定函数绑定请求数据
// 如果目标实现了Default接口，则调用Default方法设置默认值
// 最后执行传入的验证器对数据进行校验
func ReadRequest[T any](c *gin.Context, rq *T, binder Binder, validators ...Validator[T]) error {
	if err := binder(rq); err != nil {
		return errorsx.ErrBind.WithMessage(err.Error())
	}
	// 如果目标实现了Default接口，则调用Default方法设置默认值
	if defaulter, ok := any(rq).(interface{ Default() }); ok {
		defaulter.Default()
	}
	// 执行所有验证函数
	for _, validate := range validators {
		if validate == nil { // 跳过 nil 的验证器
			continue
		}
		if err := validate(c.Request.Context(), rq); err != nil {
			return err
		}
	}

	return nil
}

// 根据是否错误生成成功的响应或标准化的错误响应
func WriteResponse(c *gin.Context, data any, err error) {
	if err != nil {
		// 如果发生错误，生成错误响应
		errx := errorsx.FromError(err) // 提取错误详细信息
		c.JSON(errx.Code, ErrorResponse{
			Reason:   errx.Reason,
			Message:  errx.Message,
			Metadata: errx.MetaData,
		})
		return
	}

	// 如果没有错误，返回成功响应
	c.JSON(http.StatusOK, data)
}
