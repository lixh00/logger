# zap_logger 日志组件库使用文档

## 项目概述
`zap_logger` 是一个基于 [Zap](https://github.com/uber-go/zap) 的日志组件库，提供灵活的多输出目标支持、可定制的日志格式和级别控制，以及与 Loki 日志系统的集成能力。该库旨在简化 Go 应用程序的日志配置流程，满足不同环境下的日志收集需求。

## 核心功能
- **多输出目标**：同时支持控制台输出、文件输出和 Loki 日志系统
- **灵活配置**：通过 YAML 配置文件或代码选项动态调整日志行为
- **日志轮转**：文件日志自动轮转，支持大小限制、备份数量和压缩
- **日志编码**：支持 JSON 和 Console 两种编码格式，可按输出目标单独配置
- **级别控制**：支持细粒度的日志级别设置（Debug, Info, Warn, Error 等）
- **默认配置**：内置合理默认值，零配置快速启动

## 安装方法
```go
go get gitea.mrx.ltd/pkg/zap_logger
```

## 快速开始

### 默认配置
无需任何配置文件，直接初始化即可使用默认配置（控制台输出，Info级别，Console编码）：

```go
package main

import (
  "gitea.mrx.ltd/pkg/zap_logger"
  "go.uber.org/zap"
)

func main() {
  // 使用默认配置初始化
  if err := zap_logger.NewZapLogger(""); err != nil {
    panic(err)
  }
 

  // 使用全局日志器
  log.Info("这是一条info级别的日志")
  log.Error("这是一条error级别的日志")
}
```

### 自定义配置文件
创建 YAML 配置文件（如 `logger.yaml`）：

```yaml
logger:
  encoder: "json"       # 全局编码器 (json/console)
  level: "info"         # 全局日志级别

file:
  enable: true          # 启用文件输出,默认关闭
  filename: "app.log"   # 日志文件名,不填写则自动生成
  max_size: 10          # 单个文件最大尺寸(MB)
  max_age: 7            # 日志保留天数
  max_backups: 5        # 最大备份数量
  local_time: true      # 使用本地时间命名备份文件
  compress: true        # 压缩备份文件

console:
  enable: true          # 启用控制台输出,默认开启
  color: true           # 启用彩色输出

loki:
  enable: true          # 启用Loki输出,默认关闭
  host: "localhost"     # Loki服务地址
  port: 3100            # Loki服务端口
  source: "app"         # 日志来源标识
  service: "api"        # 服务名称
  job: "backend"        # 任务名称
  environment: "prod"   # 环境标识
```

通过配置文件初始化：

```go
if err := zap_logger.NewZapLogger("logger.yaml"); err != nil {
  panic(err)
}
```

## 配置说明

### 全局配置 (logger)
| 字段名   | 类型   | 说明                          | 可选值                 | 默认值    |
|----------|--------|-------------------------------|------------------------|-----------|
| encoder  | string | 全局日志编码器                | "json", "console"      | "console" |
| level    | string | 全局日志级别                  | "debug", "info", "warn", "error", "dpanic", "panic", "fatal" | "info"    |

### 文件输出配置 (file)
| 字段名        | 类型   | 说明                          | 默认值                          |
|---------------|--------|-------------------------------|---------------------------------|
| enable        | bool   | 是否启用文件输出              | false                           |
| encoder       | string | 文件日志编码器（覆盖全局）    | 继承 logger.encoder             |
| level         | string | 文件日志级别（覆盖全局）      | 继承 logger.level               |
| filename      | string | 日志文件路径                  | "app-{hostname}-{date}.log"     |
| max_size      | int    | 单个文件最大尺寸(MB)          | 10                              |
| max_age       | int    | 日志保留天数                  | 7                               |
| max_backups   | int    | 最大备份数量                  | 5                               |
| local_time    | bool   | 备份文件使用本地时间命名      | true                            |
| compress      | bool   | 是否压缩备份文件              | false                           |

### 控制台输出配置 (console)
| 字段名   | 类型   | 说明                          | 默认值               |
|----------|--------|-------------------------------|-------------------|
| enable   | bool   | 是否启用控制台输出            | true              |
| encoder  | string | 控制台日志编码器（覆盖全局）  | 继承 logger.encoder |
| level    | string | 控制台日志级别（覆盖全局）    | 继承 logger.level   |
| color    | bool   | 是否启用彩色输出              | false             |

### Loki输出配置 (loki)
| 字段名        | 类型   | 说明                          | 默认值                          |
|---------------|--------|-------------------------------|---------------------------------|
| enable        | bool   | 是否启用Loki输出              | false                           |
| encoder       | string | Loki日志编码器（覆盖全局）    | 继承 logger.encoder             |
| level         | string | Loki日志级别（覆盖全局）      | 继承 logger.level               |
| host          | string | Loki服务主机地址              | "localhost"                     |
| port          | int    | Loki服务端口                  | 3100                            |
| source        | string | 日志来源标识                  | ""                              |
| service       | string | 服务名称标签                  | ""                              |
| job           | string | 任务名称标签                  | ""                              |
| environment   | string | 环境标识标签                  | ""                              |

## 代码配置选项
除了配置文件外，还可以通过代码选项动态调整配置：

```go
import (
  "gitea.mrx.ltd/pkg/zap_logger"
)

func main() {
  // 使用代码选项自定义配置
  err := zap_logger.NewZapLogger("",
    zap_logger.WithEncoder(zap_logger.JsonEncoder),  // 设置全局编码器为JSON
    zap_logger.WithLevel("debug"),                  // 设置全局级别为Debug
    zap_logger.WithEnableFile(true),                // 启用文件输出
    zap_logger.WithFilename("app.log"),             // 设置日志文件名
    zap_logger.WithLokiEnable(true),                // 启用Loki输出
    zap_logger.WithLokiHost("loki.example.com"),    // 设置Loki主机
  )
  if err != nil {
    panic(err)
  }
}
```

完整选项列表：

| 函数名                     | 说明                          | 参数类型/示例                  |
|----------------------------|-------------------------------|--------------------------------|
| WithEncoder                | 设置全局编码器                | JsonEncoder/ConsoleEncoder     |
| WithLevel                  | 设置全局日志级别              | "debug", "info", "warn"等      |
| WithEnableFile             | 启用/禁用文件输出             | true/false                     |
| WithFilename               | 设置日志文件名                | "app.log"                      |
| WithFileMaxSize            | 设置文件最大尺寸(MB)          | 20                             |
| WithFileMaxAge             | 设置日志保留天数              | 15                             |
| WithFileMaxBackups         | 设置最大备份数量              | 10                             |
| WithFileLocaltime          | 备份文件使用本地时间          | true/false                     |
| WithFileCompress           | 启用/禁用备份压缩             | true/false                     |
| WithConsoleEnable          | 启用/禁用控制台输出           | true/false                     |
| WithConsoleEnableColor     | 启用/禁用控制台彩色输出       | true/false                     |
| WithLokiEnable             | 启用/禁用Loki输出             | true/false                     |
| WithLokiHost               | 设置Loki服务主机              | "loki.example.com"             |
| WithLokiPort               | 设置Loki服务端口              | 3100                           |
| WithLokiSource             | 设置日志来源标识              | "payment-service"              |
| WithLokiService            | 设置服务名称标签              | "api"                          |
| WithLokiJob                | 设置任务名称标签              | "backend"                      |
| WithLokiEnvironment        | 设置环境标识标签              | "production"                   |

## 高级使用示例

### 多输出组合
同时启用文件和Loki输出，禁用控制台：

```yaml
logger:
  encoder: "json"
  level: "info"

file:
  enable: true
  filename: "app.log"

console:
  enable: false

loki:
  enable: true
  host: "loki.internal"
  port: 3100
  service: "user-service"
  environment: "prod"
```

## 注意事项
1. 配置优先级：代码选项 > 配置文件 > 默认配置
2. 日志级别：如果全局级别高于输出特定级别，以全局级别为准
3. Loki依赖：使用Loki输出时需确保Loki服务可访问
4. 错误处理：NewZapLogger返回的错误需妥善处理，避免日志初始化失败