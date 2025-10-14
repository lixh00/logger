package log

import (
	"go.uber.org/zap"
)

func Info(msg string) {
	zap.S().Info(msg)
}

func Warn(msg string) {
	zap.S().Warn(msg)
}

func Error(msg string) {
	defer zap.S().Sync()
	zap.S().Error(msg)
}

func Fatal(msg string) {
	defer zap.S().Sync()
	zap.S().Fatal(msg)
}

func Panic(msg string) {
	defer zap.S().Sync()
	zap.S().Panic(msg)
}

func Debug(msg string) {
	zap.S().Debug(msg)
}

func Infof(format string, args ...any) {
	zap.S().Infof(format, args...)
}

func Warnf(format string, args ...any) {
	zap.S().Warnf(format, args...)
}

func Errorf(format string, args ...any) {
	defer zap.S().Sync()
	zap.S().Errorf(format, args...)
}

func Fatalf(format string, args ...any) {
	defer zap.S().Sync()
	zap.S().Fatalf(format, args...)
}

func Panicf(format string, args ...any) {
	defer zap.S().Sync()
	zap.S().Panicf(format, args...)
}

func Debugf(format string, args ...any) {
	zap.S().Debugf(format, args...)
}

func Infow(msg string, keysAndValues ...any) {
	zap.S().Infow(msg, keysAndValues...)
}

func Warnw(msg string, keysAndValues ...any) {
	zap.S().Warnw(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...any) {
	defer zap.S().Sync()
	zap.S().Errorw(msg, keysAndValues...)
}

func Fatalw(msg string, keysAndValues ...any) {
	defer zap.S().Sync()
	zap.S().Fatalw(msg, keysAndValues...)
}

func Panicw(msg string, keysAndValues ...any) {
	defer zap.S().Sync()
	zap.S().Panicw(msg, keysAndValues...)
}

func Debugw(msg string, keysAndValues ...any) {
	zap.S().Debugw(msg, keysAndValues...)
}
