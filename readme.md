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

### 环境变量
```shell
export LOG_MODE=0 # 0: dev, 1: prod
export LOG_LOKI_ENABLE=1 # 是否启用Loki 0: disable, 1: enable
export LOG_FILE_ENABLE=0 # 是否启用输出到文件 0: disable, 1: enable
export LOG_LOKI_HOST=10.0.0.31 # Loki地址
export LOG_LOKI_PORT=3100 # Loki端口
export LOG_LOKI_SOURCE_NAME=tests # Loki Source 名称
export LOG_LOKI_JOB_NAME=testj # Loki Job 名称
```