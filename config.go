package logger

import "fmt"

type mode int

var (
	Dev  mode = 0
	Prod mode = 1
)

// LogConfig 日志配置
type LogConfig struct {
	Mode       mode   `env:"LOG_MODE"`             // dev, prod
	LokiEnable bool   `env:"LOG_LOKI_ENABLE"`      // 是否启用Loki
	FileEnable bool   `env:"LOG_FILE_ENABLE"`      // 是否输出到文件
	LokiHost   string `env:"LOG_LOKI_HOST"`        // Loki地址
	LokiPort   int    `env:"LOG_LOKI_PORT"`        // Loki端口
	LokiSource string `env:"LOG_LOKI_SOURCE_NAME"` // Loki的source名称
	LokiJob    string `env:"LOG_LOKI_JOB_NAME"`    // Loki的job名称
}

func (c LogConfig) getLokiPushURL() string {
	return fmt.Sprintf("http://%v:%v/loki/api/v1/push", c.LokiHost, c.LokiPort)
}
