package log

// 设计日志包
// 设计日志接口
// 定义日志级别
// 定义日志结构体类型
// 实现创建xxxLogger实例方法以及实现Logger接口方法

// 实现结构化记录方法（Infow）
// 放在Internal/pkg是因为日志包封装了一些定制化逻辑，不适合对外暴露，所以不放在pkg
// 但日志包是项目内公共使用的，所以放在internal/pkg

import (
	"context"
	"sync"
	"time"

	"github.com/ArthurWang23/miniblog/internal/pkg/contextx"
	"github.com/ArthurWang23/miniblog/internal/pkg/known"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Debugw(msg string, kvs ...any)

	Infow(msg string, kvs ...any)

	Warnw(msg string, kvs ...any)

	Errorw(msg string, kvs ...any)

	Panicw(msg string, kvs ...any)

	Fatalw(msg string, kvs ...any)
	// Sync用于刷新日志缓冲区，确保日志被完整写入目标存储
	Sync()
}

type zapLogger struct {
	z *zap.Logger
}

// 确保*zapLogger实现了Logger的接口
var _ Logger = (*zapLogger)(nil)

var (
	mu sync.Mutex
	// 定义了全局logger
	std = New(NewOptions())
)

// 初始化全局的日志对象
// 因为会给全局变量std赋值，因此加锁
func Init(opts *Options) {
	mu.Lock()
	defer mu.Unlock()
	std = New(opts)
}

func New(opts *Options) *zapLogger {
	if opts == nil {
		opts = NewOptions()
	}

	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(opts.Level)); err != nil {
		// 指定了非法日志级别，则默认使用info
		zapLevel = zapcore.InfoLevel
	}
	// 创建encoder配置，用于控制日志的输出格式
	encoderConfig := zap.NewProductionEncoderConfig()
	// 自定义MessageKey为message,message语义更明确
	// MessageKey：日志消息字段名
	encoderConfig.MessageKey = "message"
	encoderConfig.TimeKey = "timestamp"
	// 时间戳格式
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	// 时间单位
	encoderConfig.EncodeDuration = func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendFloat64(float64(d) / float64(time.Millisecond))
	}

	cfg := &zap.Config{
		DisableCaller:     opts.DisableCaller,
		DisableStacktrace: opts.DisableStacktrace,
		Level:             zap.NewAtomicLevelAt(zapLevel),
		Encoding:          opts.Format,
		EncoderConfig:     encoderConfig,
		OutputPaths:       opts.OutputPaths,
		ErrorOutputPaths:  []string{"stderr"},
	}
	// cfg 构建*zap.Logger对象，用log包封装zap包，因此在调用栈中跳过的调用深度应该增加2
	// log包的调用栈深度为1，zap包的调用栈深度为2，因此需要增加2
	z, err := cfg.Build(zap.AddStacktrace(zapcore.PanicLevel), zap.AddCallerSkip(2))
	if err != nil {
		panic(err)
	}
	// 如何将传统log包无缝接入到结构化日志系统中？
	// 将go标准库的log重定向到zap.Logger，全局重定向，自动处理所有log包调用
	zap.RedirectStdLog(z)

	return &zapLogger{z: z}
}

func Sync() {
	std.Sync()
}

func (l *zapLogger) Sync() {
	_ = l.z.Sync()
}

// 为了便于通过log.Debugw()输出日志 实现包级别函数Debugw
// 通过调用*zapLogger类型的全局变量std的Debugw方法来输出debug级别日志
// 直接调用std的Debugw而不是std.z.Sugar().Debugw() 这样可以复用*zapLogger类型的Debugw方法现在以及未来可能的实现逻辑
func Debugw(msg string, kvs ...any) {
	std.Debugw(msg, kvs...)
}

func (l *zapLogger) Debugw(msg string, kvs ...any) {
	l.z.Sugar().Debugw(msg, kvs...)
}

func Infow(msg string, kvs ...any) {
	std.Infow(msg, kvs...)
}

func (l *zapLogger) Infow(msg string, kvs ...any) {
	l.z.Sugar().Infow(msg, kvs...)
}

func Warnw(msg string, kvs ...any) {
	std.Warnw(msg, kvs...)
}

func (l *zapLogger) Warnw(msg string, kvs ...any) {
	l.z.Sugar().Warnw(msg, kvs...)
}

func Errorw(msg string, kvs ...any) {
	std.Errorw(msg, kvs...)
}

func (l *zapLogger) Errorw(msg string, kvs ...any) {
	l.z.Sugar().Errorw(msg, kvs...)
}

func Panicw(msg string, kvs ...any) {
	std.Panicw(msg, kvs...)
}

func (l *zapLogger) Panicw(msg string, kvs ...any) {
	l.z.Sugar().Panicw(msg, kvs...)
}

func Fatalw(msg string, kvs ...any) {
	std.Fatalw(msg, kvs...)
}

func (l *zapLogger) Fatalw(msg string, kvs ...any) {
	l.z.Sugar().Fatalw(msg, kvs...)
}

// W方法，withContext简称
// 由于log包会被多个请求并发调用，因此为防止id污染，每个请求都会对log包深拷贝
func W(ctx context.Context) Logger {
	return std.W(ctx)
}

func (l *zapLogger) W(ctx context.Context) Logger {
	lc := l.clone()
	// 定义一个映射，关联context提取函数和日志字段名
	contextExtractors := map[string]func(context.Context) string{
		known.XRequestID: contextx.RequestID,
		known.XUserID:    contextx.UserID,
	}
	// 遍历映射，从context中提取值并添加到日志中
	for fieldName, extractor := range contextExtractors {
		if val := extractor(ctx); val != "" {
			lc.z = lc.z.With(zap.String(fieldName, val))
		}
	}
	return lc
}

func (l *zapLogger) clone() *zapLogger {
	newLogger := *l
	return &newLogger
}
