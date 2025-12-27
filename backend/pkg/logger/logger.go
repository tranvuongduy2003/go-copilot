package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/tranvuongduy2003/go-copilot/pkg/config"
)

type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	With(fields ...Field) Logger
	Sync() error
}

type Field = zap.Field

var (
	String   = zap.String
	Strings  = zap.Strings
	Int      = zap.Int
	Int64    = zap.Int64
	Float64  = zap.Float64
	Bool     = zap.Bool
	Err      = zap.Error
	Any      = zap.Any
	Duration = zap.Duration
	Time     = zap.Time
)

type zapLogger struct {
	zap *zap.Logger
}

func New(cfg config.LogConfig) (Logger, error) {
	level := parseLevel(cfg.Level)

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var encoder zapcore.Encoder
	if cfg.Format == "console" {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		level,
	)

	zapLog := zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	return &zapLogger{zap: zapLog}, nil
}

func NewWithSampling(cfg config.LogConfig, tickMs int, first int, thereafter int) (Logger, error) {
	level := parseLevel(cfg.Level)

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var encoder zapcore.Encoder
	if cfg.Format == "console" {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	core := zapcore.NewSamplerWithOptions(
		zapcore.NewCore(
			encoder,
			zapcore.AddSync(os.Stdout),
			level,
		),
		time.Duration(tickMs)*time.Millisecond,
		first,
		thereafter,
	)

	zapLog := zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	return &zapLogger{zap: zapLog}, nil
}

func NewProduction(cfg config.LogConfig) (Logger, error) {
	return NewWithSampling(cfg, 100, 10, 10)
}

func NewDevelopment() (Logger, error) {
	return New(config.LogConfig{
		Level:  "debug",
		Format: "console",
	})
}

func (l *zapLogger) Debug(msg string, fields ...Field) {
	l.zap.Debug(msg, fields...)
}

func (l *zapLogger) Info(msg string, fields ...Field) {
	l.zap.Info(msg, fields...)
}

func (l *zapLogger) Warn(msg string, fields ...Field) {
	l.zap.Warn(msg, fields...)
}

func (l *zapLogger) Error(msg string, fields ...Field) {
	l.zap.Error(msg, fields...)
}

func (l *zapLogger) Fatal(msg string, fields ...Field) {
	l.zap.Fatal(msg, fields...)
}

func (l *zapLogger) With(fields ...Field) Logger {
	return &zapLogger{zap: l.zap.With(fields...)}
}

func (l *zapLogger) Sync() error {
	return l.zap.Sync()
}

func parseLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}
