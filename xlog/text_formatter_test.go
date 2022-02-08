package xlog

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"testing"
	"time"
)

func TestTextLoggerFormatter_Format(t *testing.T) {
	//xtrace.InitOnlyTracingLog()
	//xtrace.InitOnlyTracingLog("test")
	InitOnlyTracingLog("test")
	opentracing.StartSpan("test")
	ctx := opentracing.ContextWithSpan(context.Background(), opentracing.StartSpan("test"))

	//log = CtxLog(nil)
	SetFormatter(&TextLoggerFormatter{})
	SetReportCaller(true)
	log := CtxLog(ctx)
	log.Info("123", 1, "qqq")
	dt, _ := time.ParseInLocation(
		"20060102",
		time.Now().AddDate(0, 0, 1).Format("20060102"),
		time.Local,
	)
	fmt.Println(dt.Format("2006-01-02"))
}
