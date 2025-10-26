package zap_logger

import (
	"os"
	"strings"
	"time"

	customencoder "code.mrx.ltd/pkg/zap_logger/encoder"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type consoleLogger struct {
	Encoder string
	Level   string
	Color   bool
}

func newConsoleLogger(consoleConf *Console) *consoleLogger {
	return &consoleLogger{
		Encoder: consoleConf.Encoder,
		Level:   consoleConf.Level,
		Color:   consoleConf.Color,
	}
}

func (c *consoleLogger) Init() zapcore.Core {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(time.Format("2006-01-02 15:04:05.000"))
	}

	writer := zapcore.AddSync(os.Stdout)
	var encoder zapcore.Encoder
	switch c.Encoder {
	case "json":
		encoder = customencoder.NewJsonEncoder(encoderConfig)
	case "console":
		if c.Color {
			encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		}
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	level, _ := zapcore.ParseLevel(strings.ToLower(c.Level))

	return zapcore.NewCore(encoder, writer, level)
}
