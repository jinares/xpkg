package xlog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc/status"
	"os"
	"reflect"
	"strings"
	"time"

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
	WithFields      = xLOG.WithFields
)

func Err(err error, l ...*logrus.Entry) *logrus.Entry {
	ml := WithFields(logrus.Fields{})
	if len(l) > 0 {
		ml = l[0]
	}
	if err == nil {
		return ml
	}
	return ml.WithFields(logrus.Fields{
		"err_string": errString(err),
		"code":       errCode(err),
	})
}
func Elapsed(st time.Time, l ...*logrus.Entry) *logrus.Entry {
	if len(l) == 0 || l[0] == nil {
		return WithFields(logrus.Fields{
			"exec_time": time.Since(st).Milliseconds(),
		})
	}
	return l[0].WithField("exec_time", time.Since(st).Milliseconds())
}
func init() {
	xLOG.SetOutput(os.Stderr)
	xLOG.SetLevel(logrus.TraceLevel)
	xLOG.SetFormatter(&TextLoggerFormatter{})
}

//String err to string
func errString(err error) string {
	if err == nil {
		return ""
	}
	p := fromXErr(err)
	return fmt.Sprintf("%d %s %s", int32(p.Code()), p.Message(), jsontostr(p.Details()))
}
func fromXErr(err error) *status.Status {
	rpcerr, isok := status.FromError(err)
	if isok {
		return rpcerr
	}
	return rpcerr
}
func errCode(err error) int64 {
	return int64(fromXErr(err).Code())
}

func jsontostr(dstr interface{}) string {
	if dstr == nil {
		return ""
	}
	switch vv := dstr.(type) {
	case string:
		return vv
	}
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.Encode(dstr)

	return strings.TrimSpace(string(buffer.Bytes()))

}
func has(key interface{}, val ...interface{}) bool {
	skeyType := reflect.TypeOf(key).String()

	for _, itemVal := range val {
		if skeyType != reflect.TypeOf(itemVal).String() {
			continue
		}
		if key == itemVal {
			return true
		}
	}
	return false
}
