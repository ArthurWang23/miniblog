package apiserver

import (
	clientv3 "go.etcd.io/etcd/client/v3"

	etcd "github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/registry"
)

// NewEtcdRegistry 基于配置创建 etcd Registrar 与 Discovery（共用一个 etcd 实例）
func NewEtcdRegistry(cfg *Config) (registry.Registrar, registry.Discovery, error) {
	eo := cfg.EtcdOptions
	if eo == nil {
		return nil, nil, nil
	}

	c := clientv3.Config{
		Endpoints:   eo.Endpoints,
		DialTimeout: eo.DialTimeout,
		Username:    eo.Username,
		Password:    eo.Password,
	}

	// TLS（如有启用）
	if eo.TLSOptions != nil && eo.TLSOptions.UseTLS {
		c.TLS = eo.TLSOptions.MustTLSConfig()
	}

	cli, err := clientv3.New(c)
	if err != nil {
		return nil, nil, err
	}

	r := etcd.New(cli)
	return r, r, nil
}