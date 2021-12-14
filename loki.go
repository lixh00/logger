package logger

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/lixh00/loki-client-go/loki"
	"github.com/prometheus/common/model"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

// Loki连接对象
var lokiClient *loki.Client

// 日志输出
type lokiWriter struct{}

// 初始化LokiCore，使日志可以推送到Loki
func initLokiCore() zapcore.Core {
	initLokiClient()
	// 日志输出到控制台和Loki
	writer := zapcore.AddSync(newLokiWriter())

	// 自定义时间输出格式
	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}

	// 格式相关的配置
	encoderConfig := zap.NewProductionEncoderConfig()
	// 修改时间戳的格式
	encoderConfig.EncodeTime = customTimeEncoder
	// 日志级别使用大写
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// 将日志级别设置为 DEBUG
	return zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), writer, zapcore.DebugLevel)
}

// 初始化LokiClient
func initLokiClient() {
	// 如果Loki配置错误，返回一个nil
	if config.LokiHost == "" || config.LokiPort < 1 {
		panic(errors.New("Loki配置错误"))
	}
	// 初始化配置
	cfg, _ := loki.NewDefaultConfig(config.getLokiPushURL())
	// 创建连接对象
	client, err := loki.NewWithLogger(cfg, log.NewNopLogger())
	if err != nil {
		panic("Loki初始化失败: " + err.Error())
	}
	lokiClient = client
}

// 实现Write接口，使lokiClient可以作为zap的扩展
func (c lokiWriter) Write(p []byte) (int, error) {
	type logInfo struct {
		Level  string `json:"level"`  // 日志级别
		Ts     string `json:"ts"`     // 格式化后的时间(在zap那边配置的)
		Caller string `json:"caller"` // 日志输出的文件名和行号
		Msg    string `json:"msg"`    // 日志内容
	}
	var li logInfo
	err := json.Unmarshal(p, &li)
	if err != nil {
		return 0, err
	}

	label := model.LabelSet{"job": model.LabelValue(config.LokiName)}
	label["source"] = model.LabelValue(config.LokiName)
	label["level"] = model.LabelValue(li.Level)
	label["caller"] = model.LabelValue(li.Caller)
	// 异步推送消息到服务器
	go func() {
		err = lokiClient.Handle(label, time.Now().Local(), li.Msg)
		if err != nil {
			fmt.Printf("日志推送到Loki失败: %v\n", err.Error())
		}
	}()

	return 0, nil
}

// NewLokiWriter
func newLokiWriter() *lokiWriter {
	return &lokiWriter{}
}
