package log

import (
	"github.com/spf13/pflag"
	"go.uber.org/zap/zapcore"
)

type Options struct {
	DisableCaller bool `json:"disable_caller,omitempty" mapstructure:"disable-caller"`

	DisableStacktrace bool `json:"disable-stacktrace,omitempty" mapstructure:"disable-stacktrace"`

	EnableColor bool `json:"enable-color,omitempty" mapstructure:"enable-color"`

	Level string `json:"level,omitempty" mapstructure:"level"`

	Format string `json:"format,omitempty" mapstructure:"format"`

	OutputPaths []string `json:"output-paths,omitempty" mapstructure:"output-paths"`
}

func NewOptions() *Options {
	return &Options{
		Level:       zapcore.InfoLevel.String(),
		Format:      "console",
		OutputPaths: []string{"stdout"},
	}
}

func (o *Options) Validate() []error {
	errs := []error{}
	return errs
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Level, "log.level", o.Level, "Minimum log output `LEVEL`,")
	fs.BoolVar(&o.DisableCaller, "log.disable-caller", o.DisableCaller, "Disable output of caller information in the log.")
	fs.BoolVar(&o.DisableStacktrace, "log.disable-stacktrace", o.DisableStacktrace, ""+
		"Disable the log to record a stack trace for all messages at or above panic level.")
	fs.BoolVar(&o.EnableColor, "log.enable-color", o.EnableColor, "Enable output ansi colors in plain format logs.")
	fs.StringVar(&o.Format, "log.format", o.Format, "Log output `FORMAT`, support plain or json format.")
	fs.StringSliceVar(&o.OutputPaths, "log.output-paths", o.OutputPaths, "Output paths of log.")
}
