package options

import "github.com/spf13/pflag"

var _ IOptions = (*ClientCertAuthenticationOptions)(nil)

type ClientCertAuthenticationOptions struct {
	ClientCA string `json:"client-ca-file" mapstructure:"client-ca-file"`
}

func NewClientCertAuthenticationOptions() *ClientCertAuthenticationOptions {
	return &ClientCertAuthenticationOptions{
		ClientCA: "",
	}
}

func (o *ClientCertAuthenticationOptions) Validate() []error {
	return []error{}
}

func (o *ClientCertAuthenticationOptions) AddFlags(fs *pflag.FlagSet, prefixes ...string) {
	fs.StringVar(&o.ClientCA, "client-ca-file", o.ClientCA, ""+
		"If set, any request presenting a client certificate signed by one of "+
		"the authorities in the client-ca-file is authenticated with an identity "+
		"corresponding to the CommonName of the client certificate.")
}
