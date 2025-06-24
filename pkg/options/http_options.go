package options

import (
	"time"

	"github.com/spf13/pflag"
)

var _ IOptions = (*HTTPOptions)(nil)

type HTTPOptions struct {
	Network string `json:"network" mapstructure:"network"`

	Addr string `json:"addr" mapstructure:"addr"`

	Timeout time.Duration `json:"timeout" mapstructure:"timeout"`
}

func NewHTTPOptions() *HTTPOptions {
	return &HTTPOptions{
		Network: "tcp",
		Addr:    "0.0.0.0:38443",
		Timeout: 30 * time.Second,
	}
}

func (o *HTTPOptions) Validate() []error {
	if o == nil {
		return nil
	}

	errors := []error{}

	if err := ValidateAddress(o.Addr); err != nil {
		errors = append(errors, err)
	}
	return errors
}

func (o *HTTPOptions) AddFlags(fs *pflag.FlagSet, prefixes ...string) {
	fs.StringVar(&o.Network, "http.network", o.Network, "Specify the network for the HTTP server.")
	fs.StringVar(&o.Addr, "http.addr", o.Addr, "Specify the HTTP server bind address and port.")
	fs.DurationVar(&o.Timeout, "http.timeout", o.Timeout, "Timeout for server connections.")
}

func (s *HTTPOptions) Complete() error {
	return nil
}
