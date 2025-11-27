package zap_logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	customencoder "gitee.ltd/lxh/logger/v2/encoder"
)

type fileLogger struct {
	Filename   string
	Encoder    string
	Level      string
	MaxSize    int
	MaxAge     int
	MaxBackups int
	LocalTime  bool
	Compress   bool
}

func newFileLogger(fileConf *File) *fileLogger {
	if fileConf.Filename == "" {
		hostname, _ := os.Hostname()
		fileConf.Filename = fmt.Sprintf("app-%s-%s.log", hostname, time.Now().Format("20060102"))
	}
	return &fileLogger{
		Filename:   fileConf.Filename,
		Encoder:    fileConf.Encoder,
		Level:      fileConf.Level,
		MaxSize:    fileConf.MaxSize,
		MaxAge:     fileConf.MaxAge,
		MaxBackups: fileConf.MaxBackups,
		LocalTime:  fileConf.LocalTime,
		Compress:   fileConf.Compress,
	}
}

func (f *fileLogger) Init() zapcore.Core {
	lumberLogger := &lumberjack.Logger{
		Filename:   f.Filename,
		MaxSize:    f.MaxSize,
		MaxAge:     f.MaxAge,
		MaxBackups: f.MaxBackups,
		LocalTime:  f.LocalTime,
		Compress:   f.Compress,
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(time.Format("2006-01-02 15:04:05.000"))
	}

	writer := zapcore.AddSync(lumberLogger)
	var encoder zapcore.Encoder
	switch f.Encoder {
	case "json":
		encoder = customencoder.NewJsonEncoder(encoderConfig)
	case "console":
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	level, _ := zapcore.ParseLevel(strings.ToLower(f.Level))

	return zapcore.NewCore(encoder, writer, level)
}
