package logger

import "fmt"

type mode int

var (
	Dev  mode = 0
	Prod mode = 1
)

// LogConfig 日志配置
type LogConfig struct {
	Mode       mode // dev, prod
	LokiEnable bool
	FileEnable bool
	LokiHost   string
	LokiPort   int
	LokiName   string // Loki的job和source名称
}

func (c LogConfig) getLokiPushURL() string {
	return fmt.Sprintf("http://%v:%v/loki/api/v1/push", c.LokiHost, c.LokiPort)
}
