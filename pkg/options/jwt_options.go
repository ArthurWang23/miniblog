package options

import (
	"fmt"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/spf13/pflag"
)

var _ IOptions = (*JWTOptions)(nil)

type JWTOptions struct {
	Key           string        `json:"key" mapstructure:"key"`
	Expire        time.Duration `json:"expire" mapstructure:"expire"`
	MaxRefresh    time.Duration `json:"max-refresh" mapstructure:"max-refresh"`
	SigningMethod string        `json:"signing-method" mapstructure:"signing-method"`
}

func NewJWTOptions() *JWTOptions {
	return &JWTOptions{
		Key:           "miniblog",
		Expire:        2 * time.Hour,
		MaxRefresh:    2 * time.Hour,
		SigningMethod: "HS512",
	}
}

func (o *JWTOptions) Validate() []error {
	var errs []error

	if !govalidator.StringLength(o.Key, "6", "32") {
		errs = append(errs, fmt.Errorf("--jwt.key must larger than 5 and litter than 33"))
	}
	return errs
}

func (o *JWTOptions) AddFlags(fs *pflag.FlagSet, prefixes ...string) {
	if fs == nil {
		return
	}
	fs.StringVar(&o.Key, "jwt.key", o.Key, "Private key used to sign jwt token.")
	fs.DurationVar(&o.Expire, "jwt.expire", o.Expire, "Expire time for jwt token.")
	fs.DurationVar(&o.MaxRefresh, "jwt.max-refresh", o.MaxRefresh, ""+
		"This field allows clients to refresh their token until MaxRefresh has passed.")
	fs.StringVar(&o.SigningMethod, "jwt.signing-method", o.SigningMethod, "Signing method for jwt token.")
}
