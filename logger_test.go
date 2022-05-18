package logger

import (
	"gitee.ltd/lxh/logger/log"
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
	log.Info("我是测试消息")
	//time.Sleep(5 * time.Second)
}
