package validation

import (
	"regexp"

	"github.com/ArthurWang23/miniblog/internal/apiserver/store"
	"github.com/ArthurWang23/miniblog/internal/pkg/errno"
	"github.com/google/wire"
)

// 参数校验
// 需要支持自定义复杂校验逻辑
// 可以复用已有的参数校验逻辑
// 灵活通用的校验方式
// 简单易维护

type Validator struct {
	// 有些复杂的验证逻辑可能要直接查询数据库
	// 可以一并注入进来
	store store.IStore
}

var (
	// 预编译正则表达式 提高性能
	lengthRegex = regexp.MustCompile(`^.{3,20}$`)                                        // 长度3-20
	validRegex  = regexp.MustCompile(`^[A-Za-z0-9_]+$`)                                  // 只包含字母、数字和下划线
	letterRegex = regexp.MustCompile(`[A-Za-z]`)                                         // 至少包含一个字母
	numberRegex = regexp.MustCompile(`\d`)                                               // 至少包含一个数字
	emailRegex  = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`) // 邮箱格式
	phoneRegex  = regexp.MustCompile(`^1[3-9]\d{9}$`)                                    // 中国手机号
)

var ProviderSet = wire.NewSet(New)

func New(store store.IStore) *Validator {
	return &Validator{store: store}
}

func isValidUsername(username string) bool {

	if !lengthRegex.MatchString(username) {
		return false
	}

	if !validRegex.MatchString(username) {
		return false
	}
	return true
}

func isValidPassword(password string) error {
	switch {
	// 检查新密码是否为空
	case password == "":
		return errno.ErrInvalidArgument.WithMessage("password cannot be empty")
	// 检查新密码的长度要求
	case len(password) < 6:
		return errno.ErrInvalidArgument.WithMessage("password must be at least 6 characters long")
	// 使用正则表达式检查是否至少包含一个字母
	case !letterRegex.MatchString(password):
		return errno.ErrInvalidArgument.WithMessage("password must contain at least one letter")
	// 使用正则表达式检查是否至少包含一个数字
	case !numberRegex.MatchString(password):
		return errno.ErrInvalidArgument.WithMessage("password must contain at least one number")
	}
	return nil
}

func isValidEmail(email string) error {
	if email == "" {
		return errno.ErrInvalidArgument.WithMessage("email cannot be empty")
	}

	if !emailRegex.MatchString(email) {
		return errno.ErrInvalidArgument.WithMessage("invalid email format")
	}
	return nil
}

func isValidPhone(phone string) error {
	if phone == "" {
		return errno.ErrInvalidArgument.WithMessage("phone cannot be empty")
	}

	if !phoneRegex.MatchString(phone) {
		return errno.ErrInvalidArgument.WithMessage("invalid phone format")
	}
	return nil
}
