package zap_logger

import (
	"errors"
	"os"

	_ "embed"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

//go:embed default_logger.yaml
var defaultConfigByte []byte

type IZapLogger interface {
	Init() zapcore.Core
}

type ZapLogger struct{}

func NewZapLogger(filePath string, opts ...Option) error {
	var fileBytes []byte
	var err error
	// 读取默认配置
	if filePath == "" {
		fileBytes = defaultConfigByte
	} else {
		fileBytes, err = os.ReadFile(filePath)
		if err != nil {
			return errors.New("read logger file error: " + err.Error())
		}
	}

	var config Config
	err = yaml.Unmarshal(fileBytes, &config)
	if err != nil {
		return errors.New("unmarshal logger file error: " + err.Error())
	}

	// 应用最新设置
	for _, opt := range opts {
		opt(&config)
	}

	var cores []zapcore.Core
	if config.File.Enable {
		config.File.Encoder = config.Logger.Encoder
		config.File.Level = config.Logger.Level
		cores = append(cores, newFileLogger(config.File).Init())
	}
	if config.Console.Enable {
		config.Console.Encoder = config.Logger.Encoder
		config.Console.Level = config.Logger.Level
		cores = append(cores, newConsoleLogger(config.Console).Init())
	}
	if config.Loki.Enable {
		config.Loki.Encoder = config.Logger.Encoder
		config.Loki.Level = config.Logger.Level
		cores = append(cores, newLokiLogger(config.Loki).Init())
	}
	// 如果一个都没开启这默认开启一个console,info级别
	if len(cores) <= 0 {
		config.Console.Encoder = "console"
		config.Logger.Level = "info"
		config.Console.Level = "info"
		cores = append(cores, newConsoleLogger(config.Console).Init())
	}

	logger := zap.New(
		zapcore.NewTee(cores...),          // 开启的日志核心
		zap.AddCaller(),                   // 启用调用者信息
		zap.AddCallerSkip(1),              // 调用者信息跳过
		zap.AddStacktrace(zap.ErrorLevel), // 开启panic日志错误堆栈收集
	)
	zap.ReplaceGlobals(logger)
	return nil
}
