### Logger
基于Zap整合的日志框架，可自由组合输出到Console、File、Loki

### Demo
```go
package main

import "gitee.ltd/lxh/logger"

func main() {
	logger.InitLogger(logger.LogConfig{Mode: logger.Dev, LokiEnable: false, FileEnable: true})
	logger.Say.Debug("芜湖")
}
```