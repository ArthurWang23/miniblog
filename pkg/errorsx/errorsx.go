// Copyright 2025 ArthurWang &lt;2826979176@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/arthurwang23/miniblog. The professional
// version of this repository is https://github.com/arthurwang23/miniblog.

package errorsx

// 实现errorsx错误包

import (
	"errors"
	"fmt"
	"net/http"

	httpstatus "github.com/go-kratos/kratos/v2/transport/http/status"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

// 自定义错误包
type ErrorX struct {
	// 错误码
	Code int `json:"code,omitempty"`
	// 错误发生的原因,通常为业务错误码，用于精确定位问题
	Reason string `json:"reason,omitempty"`
	// 简短的错误信息，通常可直接暴露给用户
	Message string `json:"message,omitempty"`
	// 存储与错误相关的额外元数据，可包含上下文或调试信息
	MetaData map[string]string `json:"meta_data,omitempty"`
}

func New(code int, reason string, format string, args ...any) *ErrorX {
	return &ErrorX{
		Code:    code,
		Reason:  reason,
		Message: fmt.Sprintf(format, args...),
	}
}

// 实现error接口中的Error方法
func (err *ErrorX) Error() string {
	return fmt.Sprintf("error: code = %d reason = %s message = %s metadata = %v", err.Code, err.Reason, err.Message, err.MetaData)
}

// 设置错误的Message字段
func (err *ErrorX) WithMessage(format string, args ...any) *ErrorX {
	err.Message = fmt.Sprintf(format, args...)
	return err
}

func (err *ErrorX) WithMetadata(md map[string]string) *ErrorX {
	err.MetaData = md
	return err
}

// 使用key-value对设置元数据
func (err *ErrorX) KV(kvs ...string) *ErrorX {
	if err.MetaData == nil {
		err.MetaData = make(map[string]string)
	}
	for i := 0; i < len(kvs); i += 2 {
		if i+1 < len(kvs) {
			err.MetaData[kvs[i]] = kvs[i+1]
		}
	}
	return err
}

// 将ErrorX转换为grpc的status.Status类型，生成grpc标准化的错误返回信息
func (err *ErrorX) GRPCStatus() *status.Status {
	details := errdetails.ErrorInfo{Reason: err.Reason, Metadata: err.MetaData}
	s, _ := status.New(httpstatus.ToGRPCCode(err.Code), err.Message).WithDetails(&details)
	return s
}

func (err *ErrorX) WithRequestID(requestID string) *ErrorX {
	return err.KV("X-Request-ID", requestID)
}

// Is判断当前错误是否与目标错误匹配
// 递归遍历错误链，并比较ErrorX实例的Code和Reason字段
// 若Code和Reason均相等，则返回true
func (err *ErrorX) Is(target error) bool {
	if errx := new(ErrorX); errors.As(target, &errx) {
		return errx.Code == err.Code && errx.Reason == err.Reason
	}
	return false
}

// 返回错误的http码
func Code(err error) int {
	if err == nil {
		return http.StatusOK
	}
	return FromError(err).Code
}

// 返回特定错误原因
func Reason(err error) string {
	if err == nil {
		return ErrInternal.Reason
	}
	return FromError(err).Reason
}

// 将通用error转换为自定义的*ErrorX
func FromError(err error) *ErrorX {
	if err == nil {
		return nil
	}
	// 检查传入的error是否已经是ErrorX
	// 若可以通过error.As转换，则直接返回
	if errx := new(ErrorX); errors.As(err, &errx) {
		return errx
	}

	// grpc的status.FromError方法尝试将error转换为grpc错误的status对象
	// 若err不能转换为grpc错误，则返回一个带有默认值的ErrorX表示是一个位置类型的错误
	gs, ok := status.FromError(err)
	if !ok {
		return New(ErrInternal.Code, ErrInternal.Reason, err.Error())
	}

	// 若err是grpc错误类型，会成功返回一个grpc status对象
	// 使用grpc状态中的错误代码和消息创建一个ErrorX
	ret := New(httpstatus.FromGRPCCode(gs.Code()), ErrInternal.Reason, gs.Message())
	// 遍历grpc错误详情中的所有附加信息（details）
	for _, detail := range gs.Details() {
		if typed, ok := detail.(*errdetails.ErrorInfo); ok {
			ret.Reason = typed.Reason
			return ret.WithMetadata(typed.Metadata)
		}
	}
	return ret
}
