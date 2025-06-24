package options

import "github.com/spf13/pflag"

var _ IOptions = (*ConsulOptions)(nil)

// options for consul client
type ConsulOptions struct {
	Addr string `json:"addr,omitempty" mapstructure:"addr"`

	Scheme string `json:"scheme,omitempty" mapstructure:"scheme"`
}

func NewConsulOptions() *ConsulOptions {
	return &ConsulOptions{
		Addr:   "127.0.0.1:8500",
		Scheme: "http",
	}
}

func (o *ConsulOptions) Validate() []error {
	return []error{}
}

func (o *ConsulOptions) AddFlags(fs *pflag.FlagSet, prefixes ...string) {
	fs.StringVar(&o.Addr, "consul-addr", o.Addr, ""+
		"Addr is the address of the consul server.")

	fs.StringVar(&o.Scheme, "consul-scheme", o.Scheme, ""+
		"Scheme is the URI scheme for the consul server.")
}
