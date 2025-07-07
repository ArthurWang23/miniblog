package options

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

var _ IOptions = (*TLSOptions)(nil)

// TLSOptions is the TLS cert info for serving secure traffic.
type TLSOptions struct {
	// 是否启用TLS
	UseTLS bool `json:"use-tls" mapstructure:"use-tls"`
	// 是否跳过证书校验，若true，则在客户端验证服务端证书时，会忽略证书的真实性或有效性，通常用于测试环境
	InsecureSkipVerify bool   `json:"insecure-skip-verify" mapstructure:"insecure-skip-verify"`
	ServerName         string `json:"server-name" mapstructure:"server-name"`
	// 设置CA
	CaCert string `json:"ca-file" mapstructure:"ca-file"`
	// 设置客户端或服务端的证书文件路径
	Cert string `json:"cert" mapstructure:"cert"`
	// 与Cert对应的私钥文件路径，用于证明证书持有者的身份以及完成TLS握手过程
	Key string `json:"key" mapstructure:"key"`
}

func NewTLSOptions() *TLSOptions {
	return &TLSOptions{}
}

func (o *TLSOptions) Validate() []error {
	errs := []error{}

	if !o.UseTLS {
		return errs
	}

	if (o.Cert != "" && o.Key == "") || (o.Cert == "" && o.Key != "") {
		errs = append(errs, fmt.Errorf("cert and key must be set together"))
	}

	return errs
}

func (o *TLSOptions) AddFlags(fs *pflag.FlagSet, prefixes ...string) {
	fs.BoolVar(&o.UseTLS, join(prefixes...)+"tls.use-tls", o.UseTLS, "Use tls transport to connect the server.")
	fs.BoolVar(&o.InsecureSkipVerify, join(prefixes...)+"tls.insecure-skip-verify", o.InsecureSkipVerify, ""+
		"Control whether a clinet verifies the server's certificate chain and host name.")
	fs.StringVar(&o.CaCert, join(prefixes...)+"tls.ca-cert", o.CaCert, "Path to ca cert for connecting to the server.")
	fs.StringVar(&o.Cert, join(prefixes...)+"tls.cert", o.Cert, "Path to cert file for connecting to the server.")
	fs.StringVar(&o.Key, join(prefixes...)+"tls.key", o.Key, "Path to key file for connecting to the server.")
}

func (o *TLSOptions) MustTLSConfig() *tls.Config {
	tlsConf, err := o.TLSConfig()
	if err != nil {
		return &tls.Config{}
	}
	return tlsConf
}

func (o *TLSOptions) TLSConfig() (*tls.Config, error) {
	if !o.UseTLS {
		return nil, nil
	}
	tlsConfig := &tls.Config{
		InsecureSkipVerify: o.InsecureSkipVerify,
	}

	if o.Cert != "" && o.Key != "" {
		var cert tls.Certificate
		cert, err := tls.LoadX509KeyPair(o.Cert, o.Key)
		if err != nil {
			return nil, fmt.Errorf("failed to loading tls certificates: %w", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	if o.CaCert != "" {
		data, err := os.ReadFile(o.CaCert)
		if err != nil {
			return nil, err
		}
		capool := x509.NewCertPool()
		for {
			var block *pem.Block
			block, _ = pem.Decode(data)
			if block == nil {
				break
			}
			cacert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, err
			}
			capool.AddCert(cacert)
		}
		tlsConfig.RootCAs = capool
	}
	return tlsConfig, nil
}
