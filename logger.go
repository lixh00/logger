package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var config LogConfig
var Say *zap.SugaredLogger

// 避免异常，在第一次调用时初始化一个只打印到控制台的logger
func init() {
	InitLogger(LogConfig{Mode: Dev, LokiEnable: false, FileEnable: false})
}

// InitLogger 初始化日志工具
func InitLogger(c LogConfig) {
	config = c
	var cores []zapcore.Core
	// 生成输出到控制台的Core
	consoleCore := initConsoleCore()
	cores = append(cores, consoleCore)
	// 生成输出到Loki的Core
	if config.LokiEnable {
		lokiCore := initLokiCore()
		cores = append(cores, lokiCore)
	}
	// 输出到文件的Core
	if config.FileEnable {
		fileCore := initFileCore()
		cores = append(cores, fileCore)
	}

	// 增加 caller 信息
	logger := zap.New(zapcore.NewTee(cores...), zap.AddCaller())
	Say = logger.Sugar()
}
