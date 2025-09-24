package options

import (
	"time"

	genericoptions "github.com/ArthurWang23/miniblog/pkg/options"
	"github.com/spf13/pflag"
)

type ServerOptions struct {
	HTTPOptions *genericoptions.HTTPOptions
	TLSOptions  *genericoptions.TLSOptions
	EtcdOptions *genericoptions.EtcdOptions
	GRPCOptions *genericoptions.GRPCOptions
	// 新增：MySQL 配置
	MySQLOptions *genericoptions.MySQLOptions
	// 新增：Kafka 配置
	KafkaOptions *genericoptions.KafkaOptions

	// 简化：仅提供 HTTP 服务
	ServerMode string

	// JWT 配置
	JWTKey     string
	Expiration time.Duration
}

func NewServerOptions() *ServerOptions {
	return &ServerOptions{
		HTTPOptions: genericoptions.NewHTTPOptions(),
		TLSOptions:  &genericoptions.TLSOptions{},
		EtcdOptions: &genericoptions.EtcdOptions{},
		GRPCOptions: genericoptions.NewGRPCOptions(),
		// 新增：默认 KafkaOptions
		KafkaOptions: genericoptions.NewKafkaOptions(),
		ServerMode:  "http", // 你也可以改为 "grpc" 或 "grpc-gateway" 在后续步骤中支持
		JWTKey:      "miniblog",
		Expiration:  time.Hour * 2,
	}
}

func (o *ServerOptions) AddFlags(fs *pflag.FlagSet) {
	// HTTP
	fs.StringVar(&o.HTTPOptions.Network, "http.network", o.HTTPOptions.Network, "HTTP network")
	fs.StringVar(&o.HTTPOptions.Addr, "http.addr", o.HTTPOptions.Addr, "HTTP listen address")
	fs.DurationVar(&o.HTTPOptions.Timeout, "http.timeout", o.HTTPOptions.Timeout, "HTTP timeout")

	// TLS（可选）
	fs.BoolVar(&o.TLSOptions.UseTLS, "tls", o.TLSOptions.UseTLS, "Enable TLS for server")
	fs.StringVar(&o.TLSOptions.CertFile, "tls.cert", o.TLSOptions.CertFile, "TLS cert file path")
	fs.StringVar(&o.TLSOptions.KeyFile, "tls.key", o.TLSOptions.KeyFile, "TLS key file path")

	// etcd（可选）
	fs.StringSliceVar(&o.EtcdOptions.Endpoints, "etcd.endpoints", o.EtcdOptions.Endpoints, "Etcd endpoints")
	fs.DurationVar(&o.EtcdOptions.DialTimeout, "etcd.dial-timeout", o.EtcdOptions.DialTimeout, "Etcd dial timeout")
	fs.StringVar(&o.EtcdOptions.Username, "etcd.username", o.EtcdOptions.Username, "Etcd username")
	fs.StringVar(&o.EtcdOptions.Password, "etcd.password", o.EtcdOptions.Password, "Etcd password")
	// TLS for etcd
	fs.BoolVar(&o.EtcdOptions.TLSOptions.UseTLS, "etcd.tls", o.EtcdOptions.TLSOptions.UseTLS, "Use TLS for etcd")
	fs.StringVar(&o.EtcdOptions.TLSOptions.CAFile, "etcd.tls.ca", o.EtcdOptions.TLSOptions.CAFile, "Etcd CA file")
	fs.StringVar(&o.EtcdOptions.TLSOptions.CertFile, "etcd.tls.cert", o.EtcdOptions.TLSOptions.CertFile, "Etcd cert file")
	fs.StringVar(&o.EtcdOptions.TLSOptions.KeyFile, "etcd.tls.key", o.EtcdOptions.TLSOptions.KeyFile, "Etcd key file")

	// JWT
	fs.StringVar(&o.JWTKey, "jwt.key", o.JWTKey, "JWT signing key")
	fs.DurationVar(&o.Expiration, "jwt.expiration", o.Expiration, "JWT expiration")
	// gRPC
	fs.StringVar(&o.GRPCOptions.Network, "grpc.network", o.GRPCOptions.Network, "gRPC network")
	fs.StringVar(&o.GRPCOptions.Addr, "grpc.addr", o.GRPCOptions.Addr, "gRPC listen address")
	fs.DurationVar(&o.GRPCOptions.Timeout, "grpc.timeout", o.GRPCOptions.Timeout, "gRPC timeout")

	// ServerMode: http | grpc | both
	fs.StringVar(&o.ServerMode, "server-mode", o.ServerMode, "Server mode: http | grpc | both")

	// 新增：Kafka flags
	o.KafkaOptions.AddFlags(fs)

	// 新增：MySQL flags
	o.MySQLOptions.AddFlags(fs)
}

func (o *ServerOptions) Validate() error {
	// 简化校验
	return nil
}

type Config struct {
	HTTPOptions *genericoptions.HTTPOptions
	TLSOptions  *genericoptions.TLSOptions
	EtcdOptions *genericoptions.EtcdOptions
	GRPCOptions *genericoptions.GRPCOptions

	// 新增：Kafka 配置
	KafkaOptions *genericoptions.KafkaOptions

	ServerMode string

	JWTKey     string
	Expiration time.Duration
}

func (o *ServerOptions) Config() *Config {
	return &Config{
		HTTPOptions: o.HTTPOptions,
		TLSOptions:  o.TLSOptions,
		EtcdOptions: o.EtcdOptions,
		GRPCOptions: o.GRPCOptions,
		// 新增：透传 Kafka
		KafkaOptions: o.KafkaOptions,
		// 新增：透传 MySQL
		MySQLOptions: o.MySQLOptions,
		ServerMode:  o.ServerMode,
		JWTKey:      o.JWTKey,
		Expiration:  o.Expiration,
	}
}