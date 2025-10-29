package encoder

import (
	"encoding/json"
	"strings"

	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

type jsonEncoder struct {
	zapcore.Encoder
	pool buffer.Pool
	zapcore.EncoderConfig
}

func NewJsonEncoder(enderConfig zapcore.EncoderConfig) *jsonEncoder {
	return &jsonEncoder{
		Encoder:       zapcore.NewJSONEncoder(enderConfig),
		pool:          buffer.NewPool(),
		EncoderConfig: enderConfig,
	}
}

// Clone 服用自带encoder的clone方法
func (e *jsonEncoder) Clone() zapcore.Encoder {
	clone := e.Encoder.Clone()
	return clone
}

// EncodeEntry 重写EncodeEntry方法
func (e *jsonEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	buf, err := e.Encoder.EncodeEntry(entry, fields)
	if err != nil {
		return nil, err
	}

	if entry.Stack != "" && e.StacktraceKey != "" {
		return e.formatStacktrace(buf)
	}

	return buf, nil
}

// formatStacktrace 自定义stacktrace格式
func (e *jsonEncoder) formatStacktrace(originalBuf *buffer.Buffer) (*buffer.Buffer, error) {
	var jsonData struct {
		Level      string `json:"level"`
		Time       string `json:"time"`
		Caller     string `json:"caller,omitempty"`
		Msg        string `json:"msg"`
		StackTrace any    `json:"stacktrace,omitempty"`
	}

	if err := json.Unmarshal(originalBuf.Bytes(), &jsonData); err != nil {
		return nil, err
	}

	// 存在就处理,反之不处理
	if jsonData.StackTrace != nil && e.StacktraceKey != "" {
		formattedStack := e.customStackFormat(jsonData.StackTrace.(string))
		jsonData.StackTrace = formattedStack
	}

	newBuf := e.pool.Get()
	encoder := json.NewEncoder(newBuf)
	encoder.SetEscapeHTML(true)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(jsonData); err != nil {
		return nil, err
	}

	return newBuf, nil
}

// customStackFormat 自定义stacktrace格式
func (e *jsonEncoder) customStackFormat(stack string) []map[string]string {
	lines := strings.Split(strings.TrimSpace(stack), "\n")
	var filteredLines []map[string]string

	// 是否为偶数行,如果不是添加空行
	if len(lines)%2 != 0 {
		lines = append(lines, "")
	}

	// 每两行一组：奇数行作为key，偶数行作为value
	for i := 0; i < len(lines); i += 2 {
		key := lines[i]
		value := ""
		if i+1 < len(lines) {
			value = lines[i+1]
		}
		// 取出换行符和制表符
		key = strings.Replace(key, "\n", "", -1)
		value = strings.Replace(value, "\t", "", -1)
		filteredLines = append(filteredLines, map[string]string{
			key: value,
		})
	}

	return filteredLines
}
