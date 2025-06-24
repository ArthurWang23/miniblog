package options

import "github.com/spf13/pflag"

type IOptions interface {
	Validate() []error

	AddFlags(fs *pflag.FlagSet, prefixes ...string)
}
