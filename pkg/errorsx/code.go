package errorsx

import "net/http"

var (
	OK = &ErrorX{Code: http.StatusOK, Message: ""}
	// 表示所有未知的服务器端错误
	ErrInternal = &ErrorX{Code: http.StatusInternalServerError, Reason: "InternalError", Message: "Internal server error"}
	// 表示资源未找到
	ErrNotFound = &ErrorX{Code: http.StatusNotFound, Reason: "NotFound", Message: "Resource not found"}
	// 请求体绑定错误
	ErrBind = &ErrorX{Code: http.StatusBadRequest, Reason: "BindError", Message: "Error occurred whine binging the request body to the struct"}
	// 参数验证失败
	ErrInvalidArgument = &ErrorX{Code: http.StatusBadRequest, Reason: "InvalidArgument", Message: "Argument verification failed"}
	// 认证失败
	ErrUnauthenticated = &ErrorX{Code: http.StatusUnauthorized, Reason: "Unauthenticated", Message: "Unauthenticated"}
	// 没有权限
	ErrPermissionDenied = &ErrorX{Code: http.StatusForbidden, Reason: "PermissionDenied", Message: "Permission denied.Access to the requested resource is forbidden."}
	// 操作失败
	ErrOperationFailed = &ErrorX{Code: http.StatusConflict, Reason: "OperationFailed", Message: "The requested operation has failed.Please try again later."}
)
