package log

import "go.uber.org/zap/zapcore"

// 定义日志配置的选项结构体
// 可以自定义日志的输出格式，级别以及其他相关配置
type Options struct {
	// DisableCaller 指定是否禁用caller信息
	// 若false 日志中会显示调用日志所在的文件名和行号
	DisableCaller bool

	// DisableStacktrace 指定是否禁用堆栈跟踪
	// 若false 日志级别为panic或更高时会打印堆栈跟踪信息
	DisableStacktrace bool

	// Level 指定日志级别
	// 可选值：debug,info,warn,error,panic,fatal
	Level string

	// Format 指定日志输出格式
	// 可选值：json,console
	Format string

	// OutputPaths指定日志输出位置
	// 默认值为标准输出，也可以指定文件路径或其他输出目标
	OutputPaths []string
}

// 返回带有默认值的Options对象
// 用于初始化日志配置选项，提供默认的日志级别、格式和输出位置
func NewOptions() *Options {
	return &Options{
		DisableCaller:     false,
		DisableStacktrace: false,
		Level:             zapcore.InfoLevel.String(),
		Format:            "console",
		OutputPaths:       []string{"stdout"},
	}
}
