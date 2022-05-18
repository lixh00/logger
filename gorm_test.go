package logger

import (
	"gitee.ltd/lxh/logger/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestGormLogger(t *testing.T) {
	dsn := "saas:saas123@tcp(10.11.0.10:3307)/saas_tenant?charset=utf8mb4&parseTime=True&loc=Local"

	engine, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: DefaultGormLogger()})
	if err != nil {
		log.Panicf("mysql connect error: %s", err.Error())
	}

	var count int64
	if err := engine.Table("t_tenant1").Count(&count).Error; err != nil {
		t.Log(err)
	}
	t.Logf("count: %d", count)
}
