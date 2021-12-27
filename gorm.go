package logger

import "strings"

// 基于Gorm的日志实现
type gormLogger struct{}

// 打印
func (gormLogger) Write(p []byte) (n int, err error) {
	str := string(p)
	// 去掉第一行
	//str = strings.Split(str, "\n")[1]
	str = strings.Join(strings.Split(str, "\n")[1:], " ")
	Say.Debug(str)
	return 0, nil
}

// NewGormLogger ...
func NewGormLogger() *gormLogger {
	return &gormLogger{}
}
