package zap_logger

import (
	"context"
	"errors"
	"fmt"
	"time"

	"code.mrx.ltd/pkg/logger/log"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DefaultGormLogger = NewGormLogger(logger.Config{
	SlowThreshold:             2 * time.Second,
	Colorful:                  false,
	IgnoreRecordNotFoundError: false,
	ParameterizedQueries:      false,
	LogLevel:                  logger.Info,
})

type gormLogger struct{ logger.Config }

func NewGormLogger(gl logger.Config) *gormLogger {
	return &gormLogger{gl}
}

func (g *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	g.LogLevel = level
	return g
}

func (g *gormLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	if g.LogLevel >= logger.Info {
		log.Infof(msg, args...)
	}
}

func (g *gormLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	if g.LogLevel >= logger.Warn {
		log.Warnf(msg, args...)
	}
}

func (g *gormLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	if g.LogLevel >= logger.Error {
		log.Errorf(msg, args...)
	}
}

func (g *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if g.LogLevel <= logger.Silent {
		return
	}

	sql, rows := fc()
	elapsed := time.Since(begin)
	msg := fmt.Sprintf("[%v] [rows:%v] %s", elapsed, rows, sql)
	if rows == -1 {
		msg = fmt.Sprintf("[%v] [rows:%v] %s", elapsed, "-", sql)
	}
	if elapsed > g.SlowThreshold && g.SlowThreshold != 0 && g.LogLevel >= logger.Warn {
		msg = fmt.Sprintf("[SLOW SQL] [%v] [rows:%v] %s", elapsed, rows, sql)
	}

	switch {
	case err != nil && g.LogLevel >= logger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !g.IgnoreRecordNotFoundError):
		msg = fmt.Sprintf("%s -> %s", err.Error(), msg)
		g.Error(ctx, msg)
	case elapsed > g.SlowThreshold && g.SlowThreshold != 0 && g.LogLevel >= logger.Warn:
		msg = fmt.Sprintf("SLOW SQL -> %s", msg)
		g.Warn(ctx, msg)
	case g.LogLevel == logger.Info:
		g.Info(ctx, msg)
	}
}
