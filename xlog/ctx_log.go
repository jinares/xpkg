package xlog

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
)

//ContextLog ContextLog
func CtxLog(ctx context.Context) *logrus.Entry {

	traceID, spanID, parentSpanID := "", "", ""
	if ctx == nil {
		return xLOG.WithContext(ctx)
	}
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		sc, isok := span.Context().(jaeger.SpanContext)
		if isok {
			traceID = sc.TraceID().String()
			spanID = sc.SpanID().String()
			parentSpanID = sc.ParentID().String()
		}

	}
	nlog := xLOG.WithFields(logrus.Fields{
		ctxTraceID:  traceID,
		ctxSpanID:   spanID,
		ctxParentID: parentSpanID,
	})
	return nlog
}
