package options

import "github.com/spf13/pflag"

var _ IOptions = (*InsecureServingOptions)(nil)

// InsecureServingOptions are for creating an unauthenticated, unauthorized, insecure port.
// No one should be using these anymore.
type InsecureServingOptions struct {
	Addr string `json:"addr" mapstructure:"addr"`
}

func NewInsecureServingOptions() *InsecureServingOptions {
	return &InsecureServingOptions{
		Addr: "127.0.0.1:8080",
	}
}

func (o *InsecureServingOptions) Validate() []error {
	var errors []error
	return errors
}

func (o *InsecureServingOptions) AddFlags(fs *pflag.FlagSet, prefixes ...string) {
	fs.StringVar(&o.Addr, "insecure.addr", o.Addr, "The address to listen on for the insecure server.")
}
