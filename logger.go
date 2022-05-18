package logger

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var config LogConfig
var Say *zap.SugaredLogger

// 避免异常，在第一次调用时初始化一个只打印到控制台的logger
func init() {
	if Say == nil {
		// 从环境变量读取配置
		var c LogConfig
		if err := env.Parse(&c); err != nil {
			fmt.Println("日志配置解析错误: " + err.Error())
			c = LogConfig{Mode: Dev, LokiEnable: false, FileEnable: false}
		}
		// 如果值错了，直接默认为Prod
		if c.Mode != Dev && c.Mode != Prod {
			c.Mode = Prod
		}
		InitLogger(c)
	}
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
	// AddCallerSkip 输出的文件名和行号是调用封装函数的位置，而不是调用日志函数的位置
	logger := zap.New(zapcore.NewTee(cores...), zap.AddCaller())
	Say = logger.Sugar()
	zap.ReplaceGlobals(logger)
}
