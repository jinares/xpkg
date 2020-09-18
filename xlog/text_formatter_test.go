package xlog

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"testing"
)

func TestTextLoggerFormatter_Format(t *testing.T) {
	//xtrace.InitOnlyTracingLog()
	//xtrace.InitOnlyTracingLog("test")
	InitFromEnv()
	opentracing.StartSpan("test")
	ctx := opentracing.ContextWithSpan(context.Background(), opentracing.StartSpan("test"))
	log := CtxLog(ctx)
	//log = CtxLog(nil)
	SetFormatter(&TextLoggerFormatter{})
	log.Info("123", 1, "qqq")
}
