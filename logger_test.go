package logger

import (
	"gitee.ltd/lxh/logger/log"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	InitLogger(LogConfig{Mode: Dev, LokiEnable: false, FileEnable: true})
	log.Debug("芜湖")
}

func TestLogger1(t *testing.T) {
	log.Info("我是测试消息")
	time.Sleep(5 * time.Second)
}

func TestLogger2(t *testing.T) {
	InitLogger(LogConfig{Mode: Dev, LokiEnable: false, FileEnable: true})
	log.Info("我是测试消息")
	//time.Sleep(5 * time.Second)
}

func TestLogger3(t *testing.T) {
	InitLogger(LogConfig{
		Mode:       Dev,
		LokiEnable: true,
		FileEnable: false,
		LokiHost:   "",
		LokiPort:   0,
		LokiSource: "test-logger",
		LokiJob:    "test-logger",
	})

	log.Info("这是info日志")
	log.Debug("这是debug日志")
	log.Warn("这是warn日志")
}
