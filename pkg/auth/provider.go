package auth

import (
	"github.com/google/wire"
	"gorm.io/gorm"
)

// 适配 wire 的 provider，允许用 []Option 组合注入到 NewAuthz.
func NewAuthzWithOptions(db *gorm.DB, opts []Option) (*Authz, error) {
	return NewAuthz(db, opts...)
}

// ProviderSet 暴露默认选项和授权器构造给 wire.
var ProviderSet = wire.NewSet(
	DefaultOptions,      // 提供 []Option
	NewAuthzWithOptions, // 使用 []Option 构造 *Authz
)
