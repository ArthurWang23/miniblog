package errno

import (
	"net/http"

	"github.com/ArthurWang23/miniblog/pkg/errorsx"
)

//预定义错误

var (
	OK = errorsx.ErrorX{Code: http.StatusOK, Message: ""}

	// ErrInternal表示所有未知的服务器端错误
	ErrInternal = errorsx.ErrInternal

	ErrNotFound = errorsx.ErrNotFound

	ErrBind = errorsx.ErrBind

	ErrInvalidArgument = errorsx.ErrInvalidArgument

	ErrUnauthenticated = errorsx.ErrUnauthenticated

	ErrPermissionDenied = errorsx.ErrPermissionDenied

	ErrOperationFailed = errorsx.ErrOperationFailed
	// 页面没找到
	ErrPageNotFound = errorsx.ErrorX{Code: http.StatusNotFound, Reason: "NotFound.PageNotFound", Message: "Page not found"}

	// ErrSignToken 表示签发 JWT Token 时出错.
	ErrSignToken = &errorsx.ErrorX{Code: http.StatusUnauthorized, Reason: "Unauthenticated.SignToken", Message: "Error occurred while signing the JSON web token."}

	// ErrTokenInvalid 表示 JWT Token 格式无效.
	ErrTokenInvalid = &errorsx.ErrorX{Code: http.StatusUnauthorized, Reason: "Unauthenticated.TokenInvalid", Message: "Token was invalid."}

	// ErrDBRead 表示数据库读取失败.
	ErrDBRead = &errorsx.ErrorX{Code: http.StatusInternalServerError, Reason: "InternalError.DBRead", Message: "Database read failure."}

	// ErrDBWrite 表示数据库写入失败.
	ErrDBWrite = &errorsx.ErrorX{Code: http.StatusInternalServerError, Reason: "InternalError.DBWrite", Message: "Database write failure."}

	// ErrAddRole 表示在添加角色时发生错误.
	ErrAddRole = &errorsx.ErrorX{Code: http.StatusInternalServerError, Reason: "InternalError.AddRole", Message: "Error occurred while adding the role."}

	// ErrRemoveRole 表示在删除角色时发生错误.
	ErrRemoveRole = &errorsx.ErrorX{Code: http.StatusInternalServerError, Reason: "InternalError.RemoveRole", Message: "Error occurred while removing the role."}
)
