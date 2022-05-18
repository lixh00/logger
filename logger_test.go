package logger

import (
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	InitLogger(LogConfig{Mode: Dev, LokiEnable: false, FileEnable: true})
	Say.Debug("芜湖")
}

func TestLogger1(t *testing.T) {
	Say.Info("我是测试消息")
	time.Sleep(5 * time.Second)
}

func TestLogger2(t *testing.T) {
	InitLogger(LogConfig{Mode: Dev, LokiEnable: false, FileEnable: true})
	zap.S().Info("我是测试消息")
	time.Sleep(5 * time.Second)
}
