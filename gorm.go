package logger

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	gl "gorm.io/gorm/logger"
	"time"
)

var gormZap *zap.SugaredLogger

// 基于Gorm的日志实现
type gormLogger struct {
	gl.Config
}

// LogMode 实现LogMode接口
func (l *gormLogger) LogMode(level gl.LogLevel) gl.Interface {
	nl := *l
	nl.LogLevel = level
	return &nl
}

// Info 实现Info接口
func (l gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gl.Info {
		//	// 去掉第一行
		//	msg = strings.Join(strings.Split(msg, "\n")[1:], " ")
		//	gormZap.Info(msg)
		//
		//	l.Printf(msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Warn 实现Warn接口
func (l gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gl.Warn {
		//
	}
}

// 实现Error接口
func (l gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gl.Error {
		//
	}
}

// Trace 实现Trace接口
func (l gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gl.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	msg := fmt.Sprintf("[%v] [rows:%v] %s", elapsed.String(), rows, sql)
	if rows == -1 {
		msg = fmt.Sprintf("[%s] [-] %s", elapsed.String(), sql)
	}

	switch {
	case err != nil && l.LogLevel >= gl.Error && (!errors.Is(err, gl.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		gormZap.Errorf("%s -> %s", err.Error(), sql)
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gl.Warn:
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		gormZap.Warnf("%v -> %v", slowLog, sql)
	case l.LogLevel == gl.Info:
		gormZap.Info(msg)
	}
}

// NewGormLoggerWithConfig ...
func NewGormLoggerWithConfig(config gl.Config) gl.Interface {
	return &gormLogger{config}
}

// DefaultGormLogger 默认的日志实现
func DefaultGormLogger() gl.Interface {
	// 默认日志级别为Info，如果是生产环境，就是Error
	logLevel := gl.Info
	if config.Mode == Prod {
		logLevel = gl.Error
	}
	return &gormLogger{gl.Config{
		SlowThreshold:             time.Second, // Slow SQL threshold
		IgnoreRecordNotFoundError: false,       // 忽略没找到结果的错误
		LogLevel:                  logLevel,    // Log level
		Colorful:                  false,       // Disable color
	}}
}
