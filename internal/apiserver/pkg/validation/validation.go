package validation

import (
	"regexp"

	"github.com/ArthurWang23/miniblog/internal/apiserver/store"
	"github.com/ArthurWang23/miniblog/internal/pkg/errno"
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

func New(store store.IStore) *Validator {
	return &Validator{store: store}
}

func isValidUsername(username string) bool {
	var (
		lengthRegex = `^.{3,20}$`       // 长度3-20
		validRegex  = `^[A-Za-z0-9_]+$` // 只包含字母、数字和下划线
	)

	if matched, _ := regexp.MatchString(lengthRegex, username); !matched {
		return false
	}

	if matched, _ := regexp.MatchString(validRegex, username); !matched {
		return false
	}
	return true
}

func isValidPassword(password string) error {
	if password == "" {
		return errno.ErrInvalidArgument.WithMessage("password cannot be empty")
	}
	if len(password) < 6 {
		return errno.ErrInvalidArgument.WithMessage("password must be at least 6 characters long")
	}

	letterPattern := regexp.MustCompile(`[A-Za-z]`)
	if !letterPattern.MatchString(password) {
		return errno.ErrInvalidArgument.WithMessage("password must contain at least one letter")
	}

	numberPattern := regexp.MustCompile(`\d`)
	if !numberPattern.MatchString(password) {
		return errno.ErrInvalidArgument.WithMessage("password must contain at least one number")
	}

	return nil
}

func isValidEmail(email string) error {
	if email == "" {
		return errno.ErrInvalidArgument.WithMessage("email cannot be empty")
	}

	emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailPattern.MatchString(email) {
		return errno.ErrInvalidArgument.WithMessage("invalid email format")
	}
	return nil
}

func isValidPhone(phone string) error {
	if phone == "" {
		return errno.ErrInvalidArgument.WithMessage("phone cannot be empty")
	}

	phonePattern := regexp.MustCompile(`^1[3-9]\d{9}$`)
	if !phonePattern.MatchString(phone) {
		return errno.ErrInvalidArgument.WithMessage("invalid phone format")
	}
	return nil
}
