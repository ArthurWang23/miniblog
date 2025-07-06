// 授权

package auth

import (
	"time"

	casbin "github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	adapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

const (
	// 默认的Casbin访问控制模型
	defaultAclModel = `[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act, eft
	
[role_definition]
g = _, _

[policy_effect]
e = !some(where (p.eft == deny))

[matchers]
m = g(r.sub,p.sub) && keyMatch(r.obj, p.obj) && r.act == p.act`
)

// 授权器
type Authz struct {
	// 使用Casbin同步授权器
	*casbin.SyncedEnforcer
}

// 函数选项，定义NewAuthz的行为
type Option func(*authzConfig)

type authzConfig struct {
	aclModel           string        // Casbin的模型字符串
	autoLoadPolicyTime time.Duration // 自动加载策略的时间间隔
}

func defaultAuthzConfig() *authzConfig {
	return &authzConfig{
		aclModel:           defaultAclModel,
		autoLoadPolicyTime: 5 * time.Second,
	}
}

func DefaultOptions() []Option {
	return []Option{
		WithAclModel(defaultAclModel),
		WithAutoLoadPolicyTime(10 * time.Second),
	}
}

func WithAclModel(model string) Option {
	return func(cfg *authzConfig) {
		cfg.aclModel = model
	}
}

func WithAutoLoadPolicyTime(interval time.Duration) Option {
	return func(cfg *authzConfig) {
		cfg.autoLoadPolicyTime = interval
	}
}

func NewAuthz(db *gorm.DB, opts ...Option) (*Authz, error) {
	cfg := defaultAuthzConfig()
	for _, opt := range opts {
		opt(cfg)

	}

	// 初始化Gorm适配器并用于Casbin授权器
	adapter, err := adapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}
	m, _ := model.NewModelFromString(cfg.aclModel)
	enforcer, err := casbin.NewSyncedEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}
	// 从数据库加载策略
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, err
	}

	enforcer.StartAutoLoadPolicy(cfg.autoLoadPolicyTime)

	return &Authz{enforcer}, nil
}

func (a *Authz) Authorize(sub, obj, act string) (bool, error) {
	return a.Enforce(sub, obj, act)
}
