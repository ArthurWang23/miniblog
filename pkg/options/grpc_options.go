package options

import (
	"time"

	"github.com/spf13/pflag"
)

var _ IOptions = (*GRPCOptions)(nil)

type GRPCOptions struct {
	Network string `json:"network" mapstructure:"network"`

	Addr string `json:"addr" mapstructure:"addr"`

	Timeout time.Duration `json:"timeout" mapstructure:"timeout"`
}

func NewGRPCOptions() *GRPCOptions {
	return &GRPCOptions{
		Network: "tcp",
		Addr:    "0.0.0.0:39090",
		Timeout: 30 * time.Second,
	}
}

func (o *GRPCOptions) Validate() []error {
	var errors []error

	if err := ValidateAddress(o.Addr); err != nil {
		errors = append(errors, err)
	}

	return errors
}

func (o *GRPCOptions) AddFlags(fs *pflag.FlagSet, prefixes ...string) {
	fs.StringVar(&o.Network, "grpc.network", o.Network, "Specify the network for the gRPC server.")
	fs.StringVar(&o.Addr, "grpc.addr", o.Addr, "Specify the gRPC server bind address and port.")
	fs.DurationVar(&o.Timeout, "grpc.timeout", o.Timeout, "Timeout for server connections.")
}
