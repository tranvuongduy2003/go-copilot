package logger

import (
	"sync"

	"github.com/tranvuongduy2003/go-copilot/pkg/config"
)

var (
	globalLogger Logger
	once         sync.Once
)

func Init(cfg config.LogConfig, production bool) error {
	var err error
	once.Do(func() {
		if production {
			globalLogger, err = NewProduction(cfg)
		} else {
			globalLogger, err = New(cfg)
		}
	})
	return err
}

func L() Logger {
	if globalLogger == nil {
		globalLogger, _ = NewDevelopment()
	}
	return globalLogger
}

func SetLogger(l Logger) {
	globalLogger = l
}

func Debug(msg string, fields ...Field) {
	L().Debug(msg, fields...)
}

func Info(msg string, fields ...Field) {
	L().Info(msg, fields...)
}

func Warn(msg string, fields ...Field) {
	L().Warn(msg, fields...)
}

func Error(msg string, fields ...Field) {
	L().Error(msg, fields...)
}

func Fatal(msg string, fields ...Field) {
	L().Fatal(msg, fields...)
}

func With(fields ...Field) Logger {
	return L().With(fields...)
}

func Sync() error {
	if globalLogger != nil {
		return globalLogger.Sync()
	}
	return nil
}
