package zap_logger

import (
	"testing"

	"gitee.ltd/lxh/logger/v2/log"
)

func TestNewZapLogger(t *testing.T) {
	err := NewZapLogger("", WithEncoder(JsonEncoder))
	if err != nil {
		t.Fatal(err)
		return
	}
	
	log.Info("hahahahah")
	log.Error("123321")
	log.Panic("this is panic")

}
