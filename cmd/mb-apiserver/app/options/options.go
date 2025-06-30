// Copyright 2025 ArthurWang &lt;2826979176@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/arthurwang23/miniblog. The professional
// version of this repository is https://github.com/arthurwang23/miniblog.

// 通过ServerOptions结构体定义需要的配置项
// 三种配置方式：默认配置、命令行选项设置配置、通过配置文件设置
// 首先创建一个默认的配置，之后分别通过命令行选项和配置文件2种方式
// 覆盖指定的默认配置项
// 命令行选项只添加核心、必要的配置 如--config
// 其他配置通过配置文件统一配置

package options

import (
	"errors"
	"fmt"
	"time"

	genericoptions "github.com/ArthurWang23/miniblog/pkg/options"
	stringsutil "github.com/onexstack/onexstack/pkg/util/strings"
	"github.com/spf13/pflag"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/ArthurWang23/miniblog/internal/apiserver"
)

// 定义支持的服务器模式集合
// 添加grpc服务配置
var availableServerModes = sets.New(
	apiserver.GinServerMode,
	apiserver.GRPCServerMode,
	apiserver.GRPCGatewayServerMode,
)

// mapstructure 标签用于将配置文件中的配置项与go结构体字段进行映射  在调用viper.Unmarshal时会将配置项的值赋值给对应的结构体字段
type ServerOptions struct {
	// ServerMode 定义了服务器模式，可选值为grpc、Gin HTTP、HTTP Reverse Proxy
	ServerMode string `json:"server-mode" mapstructure:"server-mode"`
	// JWTKey定义JWT密钥
	JWTKey string `json:"jwt-key" mapstructure:"jwt-key"`
	// Expiration定义JWT Token过期时间
	Expiration time.Duration `json:"expiration" mapstructure:"expiration"`

	// GRPCOptions包含grpc配置选项
	GRPCOptions *genericoptions.GRPCOptions `json:"grpc" mapstructure:"grpc"`

	HTTPOptions *genericoptions.HTTPOptions `json:"http" mapstructure:"http"`

	MySQLOptions *genericoptions.MySQLOptions `json:"mysql" mapstructure:"mysql"`
}

// 创建ServerOptions的默认配置
func NewServerOptions() *ServerOptions {
	opts := &ServerOptions{
		ServerMode:   apiserver.GRPCGatewayServerMode,
		JWTKey:       "Rtg8BPKNEf2mB4mgvKONGPZZQSaJWNLijxR42qRgq0iBb5",
		Expiration:   2 * time.Hour,
		GRPCOptions:  genericoptions.NewGRPCOptions(),
		HTTPOptions:  genericoptions.NewHTTPOptions(),
		MySQLOptions: genericoptions.NewMySQLOptions(),
	}
	opts.GRPCOptions.Addr = ":6666"
	opts.HTTPOptions.Addr = ":5555"
	return opts
}

// 通过pflag从命令行解析选项   AddFlags将ServerOptions选项绑定到命令行
func (o *ServerOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.ServerMode, "server-mode", o.ServerMode, fmt.Sprintf("Server mode,available options: %v", availableServerModes.UnsortedList()))
	fs.StringVar(&o.JWTKey, "jwt-key", o.JWTKey, "JWT signing key. Must be at least 6 characters long.")
	fs.DurationVar(&o.Expiration, "expiration", o.Expiration, "The expiration duration of JWT tokens.")
	o.GRPCOptions.AddFlags(fs)
	o.HTTPOptions.AddFlags(fs)
	o.MySQLOptions.AddFlags(fs)
}

// Validate校验ServerOptions中的选项是否合法
func (o *ServerOptions) Validate() error {
	errs := []error{}
	if !availableServerModes.Has(o.ServerMode) {
		errs = append(errs, fmt.Errorf("invalid server mode: must be one of %v", availableServerModes.UnsortedList()))
	}

	if len(o.JWTKey) < 6 {
		errs = append(errs, errors.New("JWT key must be at least 6 characters long"))
	}

	if stringsutil.StringIn(o.ServerMode, []string{apiserver.GRPCServerMode, apiserver.GRPCGatewayServerMode}) {
		errs = append(errs, o.GRPCOptions.Validate()...)
	}
	// 聚合为一个错误 用的k8s生态中的一个包
	return utilerrors.NewAggregate(errs)
}

// New:运行时配置是基于初始化配置创建的，在ServerOptions中添加Config方法创建运行时配置
// 注意：导入了运行时代码包，控制面依赖数据面，要避免反向导入循环依赖
func (o *ServerOptions) Config() (*apiserver.Config, error) {
	return &apiserver.Config{
		ServerMode:   o.ServerMode,
		JWTKey:       o.JWTKey,
		Expiration:   o.Expiration,
		GRPCOptions:  o.GRPCOptions,
		HTTPOptions:  o.HTTPOptions,
		MySQLOptions: o.MySQLOptions,
	}, nil
}
