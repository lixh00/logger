package logger

import (
	"gitee.ltd/lxh/logger/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gl "gorm.io/gorm/logger"
	"testing"
	"time"
)

func TestGormLogger(t *testing.T) {
	dsn := "saas:saas123@tcp(10.11.0.10:3307)/saas_tenant?charset=utf8mb4&parseTime=True&loc=Local"

	engine, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: DefaultGormLogger()})
	if err != nil {
		log.Panicf("mysql connect error: %s", err.Error())
	}

	var count int64
	if err := engine.Table("t_tenant").Count(&count).Error; err != nil {
		t.Log(err)
	}
	t.Logf("count: %d", count)
}

func TestGormLoggerWithConfig(t *testing.T) {
	dsn := "saas:saas123@tcp(10.11.0.10:3307)/saas_tenant?charset=utf8mb4&parseTime=True&loc=Local"

	engine, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: NewGormLoggerWithConfig(gl.Config{
		SlowThreshold:             time.Second, // Slow SQL threshold
		IgnoreRecordNotFoundError: false,       // 忽略没找到结果的错误
		LogLevel:                  gl.Warn,     // Log level
		Colorful:                  false,       // Disable color
	})})
	if err != nil {
		log.Panicf("mysql connect error: %s", err.Error())
	}

	var count int64
	if err := engine.Table("t_tenant1").Count(&count).Error; err != nil {
		t.Log(err)
	}
	t.Logf("count: %d", count)
}
