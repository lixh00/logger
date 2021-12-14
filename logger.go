package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var config LogConfig
var Say *zap.SugaredLogger

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
