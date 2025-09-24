package registry

import (
	clientv3 "go.etcd.io/etcd/client/v3"

	etcd "github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/registry"

	genericoptions "github.com/ArthurWang23/miniblog/pkg/options"
)

// NewEtcdRegistryWithOptions 基于通用 EtcdOptions 创建 Registrar 与 Discovery
func NewEtcdRegistryWithOptions(eo *genericoptions.EtcdOptions) (registry.Registrar, registry.Discovery, error) {
	if eo == nil || len(eo.Endpoints) == 0 {
		return nil, nil, nil
	}

	cfg := clientv3.Config{
		Endpoints:   eo.Endpoints,
		DialTimeout: eo.DialTimeout,
		Username:    eo.Username,
		Password:    eo.Password,
	}

	// TLS（如启用）
	if eo.TLSOptions != nil && eo.TLSOptions.UseTLS {
		cfg.TLS = eo.TLSOptions.MustTLSConfig()
	}

	cli, err := clientv3.New(cfg)
	if err != nil {
		return nil, nil, err
	}

	r := etcd.New(cli)
	return r, r, nil
}