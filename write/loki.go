package write

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/grafana/loki-client-go/loki"
	"github.com/prometheus/common/model"
)

type LokiWriter struct {
	lokiClient  *loki.Client
	encoder     string
	job         string
	source      string
	environment string
}

func NewLokiWriter(lokiClient *loki.Client, encoder, job, source, environment string) *LokiWriter {
	return &LokiWriter{
		lokiClient:  lokiClient,
		encoder:     encoder,
		job:         job,
		source:      source,
		environment: environment,
	}
}

type logInfo struct {
	Level      string `json:"level"`                // 日志级别
	Ts         string `json:"time"`                 // 格式化后的时间(在zap那边配置的)
	Caller     string `json:"caller,omitempty"`     // 日志输出的文件名和行号
	Msg        string `json:"msg"`                  // 日志内容
	Stacktrace any    `json:"stacktrace,omitempty"` // 错误的堆栈信息
}

func (lw *LokiWriter) Write(p []byte) (int, error) {
	var li logInfo
	switch lw.encoder {
	case "json":
		err := json.Unmarshal(p, &li)
		if err != nil {
			return 0, err
		}
	case "console":
		li = *lw.parse(string(p))
	}

	label := model.LabelSet{"job": model.LabelValue(lw.job)}
	label["source"] = model.LabelValue(lw.source)
	label["level"] = model.LabelValue(strings.ToUpper(li.Level))
	label["caller"] = model.LabelValue(li.Caller)
	label["environment"] = model.LabelValue(lw.environment)

	t, e := time.ParseInLocation("2006-01-02 15:04:05.000", li.Ts, time.Local)
	if e != nil {
		t = time.Now().Local()
	}

	var msg string

	switch lw.encoder {
	case "json":
		jm, err := json.Marshal(li)
		if err != nil {
			return 0, err
		}

		msg = string(jm)
	case "console":
		builder := strings.Builder{}
		builder.WriteString(li.Level)
		builder.WriteString("     ")
		if li.Caller != "" {
			builder.WriteString(li.Caller)
			builder.WriteString("     ")
		}
		builder.WriteString(li.Msg)
		builder.WriteString("     ")
		if li.Stacktrace != nil {
			builder.WriteString(li.Stacktrace.(string))
		}

		msg = builder.String()
	}

	if err := lw.lokiClient.Handle(label, t, msg); err != nil {
		fmt.Printf("日志推送到Loki失败: %v\n", err.Error())
	}

	return 0, nil
}

func (lw *LokiWriter) parse(logText string) *logInfo {
	scanner := bufio.NewScanner(strings.NewReader(logText))

	if !scanner.Scan() {
		return nil // 空日志
	}

	firstLine := strings.TrimSpace(scanner.Text())
	parts := strings.Split(firstLine, "\t")

	if len(parts) < 3 {
		return nil // 格式错误
	}

	entry := &logInfo{
		Ts:    parts[0],
		Level: strings.ToUpper(parts[1]),
	}

	// 判断字段结构
	switch len(parts) {
	case 3:
		entry.Msg = parts[2]
	case 4:
		entry.Caller = parts[2]
		entry.Msg = parts[3]
	default:
		entry.Caller = parts[2]
		entry.Msg = strings.Join(parts[3:], " ")
	}

	// 收集堆栈信息
	var stackLines []string
	for scanner.Scan() {
		if scanner.Text() != "" {
			stackLines = append(stackLines, scanner.Text())
		}
	}

	if len(stackLines) > 0 {
		entry.Stacktrace = strings.Join(stackLines, "\n")
	}

	return entry
}
