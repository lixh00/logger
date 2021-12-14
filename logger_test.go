package logger

import "testing"

func TestLogger(t *testing.T) {
	InitLogger(LogConfig{Mode: Dev, LokiEnable: false, FileEnable: true})
	Say.Debug("芜湖")
}
