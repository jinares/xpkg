package xlog

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	"google.golang.org/grpc/metadata"
	"strings"
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
func Logger() *logrus.Logger {
	return xLOG
}

func TraceID(ctx context.Context) string {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		sc, isok := span.Context().(jaeger.SpanContext)
		if isok {
			return sc.TraceID().String()
		}
	}
	return RequestID(ctx)

}
func RequestID(ctx context.Context) string {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	} else {
		md = md.Copy()
	}
	mr := MDReaderWriter{md}
	traceid := ""
	mr.ForeachKey(func(key, val string) error {
		if strings.EqualFold("x_request_id", key) {
			traceid = val
			return nil
		}
		return nil
	})
	if traceid == "" {
		val := ctx.Value("x_request_id")
		switch vv := val.(type) {
		case string:
			if vv != "" {
				return vv
			}
		}
	}
	return traceid
}
