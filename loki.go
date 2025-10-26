package zap_logger

import (
	"fmt"
	"strings"
	"time"

	customencoder "code.mrx.ltd/pkg/zap_logger/encoder"
	"code.mrx.ltd/pkg/zap_logger/write"
	"github.com/grafana/loki-client-go/loki"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LokiLogger struct {
	Encoder     string `json:"encoder"`
	Level       string `json:"level"`
	Enable      bool   `yaml:"enable"`
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	Source      string `yaml:"source"`
	Service     string `yaml:"service"`
	Job         string `yaml:"job"`
	Environment string `yaml:"environment"`

	lokiClient *loki.Client
}

func newLokiLogger(lokiConf *Loki) *LokiLogger {
	lokiConfig := &LokiLogger{
		Enable:      lokiConf.Enable,
		Encoder:     lokiConf.Encoder,
		Level:       lokiConf.Level,
		Host:        lokiConf.Host,
		Port:        lokiConf.Port,
		Source:      lokiConf.Source,
		Service:     lokiConf.Service,
		Job:         lokiConf.Job,
		Environment: lokiConf.Environment,
	}

	lc, err := loki.NewWithDefault(fmt.Sprintf("%v:%v/loki/api/v1/push", lokiConfig.Host, lokiConfig.Port))
	if err != nil {
		return nil
	}

	lokiConfig.lokiClient = lc
	return lokiConfig
}

func (l *LokiLogger) Init() zapcore.Core {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(time.Format("2006-01-02 15:04:05.000"))
	}

	writer := zapcore.AddSync(write.NewLokiWriter(l.lokiClient, l.Encoder, l.Job, l.Source, l.Environment))
	var encoder zapcore.Encoder
	switch l.Encoder {
	case "json":
		encoder = customencoder.NewJsonEncoder(encoderConfig)
	case "console":
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	level, _ := zapcore.ParseLevel(strings.ToLower(l.Level))

	return zapcore.NewCore(encoder, writer, level)
}
