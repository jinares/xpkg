package xlog

import (
	"os"

	"github.com/sirupsen/logrus"
)

const (
	envServiceName = "JAEGER_SERVICE_NAME"

	ctxTraceID  = "ctx.traceid"
	ctxSpanID   = "ctx.spanid"
	ctxParentID = "ctx.parentid"
)

var (
	xLOG            = logrus.New()
	SetFormatter    = xLOG.SetFormatter
	SetOutput       = xLOG.SetOutput
	SetLevel        = xLOG.SetLevel
	SetReportCaller = xLOG.SetReportCaller
	Info            = xLOG.Info
	Infof           = xLOG.Infof
	Error           = xLOG.Error
	Errorf          = xLOG.Errorf
	Trace           = xLOG.Trace
	Tracef          = xLOG.Tracef
)

func init() {
	xLOG.SetOutput(os.Stderr)
	xLOG.SetLevel(logrus.TraceLevel)
	xLOG.SetFormatter(&TextLoggerFormatter{})
}
