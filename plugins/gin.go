package plugins

import (
	"encoding/json"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"gitee.ltd/lxh/logger/v2/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap/zapcore"
)

type LogParam struct {
	Time       string
	Method     string
	Path       string
	StatusCode int
	ClientIP   string
	UserAgent  string
	Latency    string
	Header     string
	Errors     []string
}

type FormatFunc func(logParam LogParam) string

type Skipper func(c *gin.Context) bool

var defaultConfig = &GinLogConfig{
	TimeFormat:      "2006-01-02 15:04:05.000",
	SkipPaths:       nil,
	SkipPathRegexps: nil,
	DefaultLevel:    zapcore.InfoLevel,
	Skipper:         nil,
	FormatFunc:      defaultFormatFunc,
	HideKeys:        []string{"header"},
}

var defaultFormatFunc FormatFunc = func(logParam LogParam) string {
	builder := strings.Builder{}
	builder.WriteString("[GIN]  ")
	builder.WriteString("[" + logParam.Time + "]")
	builder.WriteString("  ")
	builder.WriteString("[" + strconv.Itoa(logParam.StatusCode) + "]")
	builder.WriteString("  ")
	builder.WriteString("[" + logParam.Latency + "]")
	builder.WriteString("  ")
	builder.WriteString("[" + logParam.ClientIP + "]")
	builder.WriteString("  ")
	builder.WriteString("[" + logParam.Method + "]")
	builder.WriteString("  ")
	builder.WriteString("[" + logParam.Path + "]")
	builder.WriteString("  ")
	builder.WriteString("[" + logParam.UserAgent + "]")

	if len(logParam.Errors) > 0 {
		builder.WriteString(" | [Errors] -> ")
		for _, err := range logParam.Errors {
			builder.WriteString(" [" + err + "]")

		}

	}
	return builder.String()
}

type GinLogConfig struct {
	TimeFormat      string   // custom time format default: 2006-01-02 15:04:05.000
	SkipPaths       []string // skip path
	SkipPathRegexps []*regexp.Regexp
	DefaultLevel    zapcore.Level
	Skipper         Skipper
	FormatFunc      FormatFunc
	HideKeys        []string // current only use 'header' 'errors'
}

func GinZap() gin.HandlerFunc {
	return GinZapWithConfig(defaultConfig)
}

func GinZapWithFormat(formatFunc FormatFunc) gin.HandlerFunc {
	conf := defaultConfig
	conf.FormatFunc = formatFunc
	return GinZapWithConfig(conf)
}

func GinZapWithConfig(conf *GinLogConfig) gin.HandlerFunc {
	skipPaths := make(map[string]bool, len(conf.SkipPaths))
	for _, path := range conf.SkipPaths {
		skipPaths[path] = true
	}
	if conf.TimeFormat == "" {
		conf.TimeFormat = "2006-01-02 15:04:05.000"
	}

	if conf.FormatFunc == nil {
		conf.FormatFunc = defaultFormatFunc
	}

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.RequestURI
		c.Next()
		track := true

		if _, ok := skipPaths[path]; ok || (conf.Skipper != nil && conf.Skipper(c)) {
			track = false
		}

		if track && len(conf.SkipPathRegexps) > 0 {
			for _, reg := range conf.SkipPathRegexps {
				if !reg.MatchString(path) {
					continue
				}
				track = false
				break
			}
		}

		if track {
			latency := time.Since(start)
			responseCode := c.Writer.Status()
			method := c.Request.Method
			clientIP := c.ClientIP()
			userAgent := c.Request.UserAgent()

			logParam := LogParam{
				Time:       time.Now().Format(conf.TimeFormat),
				Method:     method,
				Path:       path,
				StatusCode: responseCode,
				ClientIP:   clientIP,
				UserAgent:  userAgent,
				Latency:    latency.String(),
			}

			if !slices.Contains(conf.HideKeys, "header") {
				headerJson, _ := json.Marshal(c.Request.Header)
				logParam.Header = string(headerJson)
			}
			if !slices.Contains(conf.HideKeys, "error") {
				var errs []string
				for _, err := range c.Errors.Errors() {
					errs = append(errs, err)
				}
				logParam.Errors = errs
			}

			if len(c.Errors.Errors()) > 0 {
				log.Error(conf.FormatFunc(logParam))
			} else {
				log.Info(conf.FormatFunc(logParam))
			}
		}
	}
}
