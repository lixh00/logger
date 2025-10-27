package zap_logger

import (
	"testing"

	"code.mrx.ltd/pkg/zap_logger/log"
)

func TestNewZapLogger(t *testing.T) {
	err := NewZapLogger("", WithEncoder(ConsoleEncoder))
	if err != nil {
		t.Fatal(err)
		return
	}

	log.Info("hahahahah")
	log.Error("123321")
	log.Panic("this is panic")

}
