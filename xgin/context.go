qpackage xgin

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jinares/xpkg/xlog"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/sirupsen/logrus"
	"strings"
)

func GinInterceptorOpenTracing(tr opentracing.Tracer) gin.HandlerFunc {
	return func(c *gin.Context) {
		if tr == nil {
			tr = opentracing.GlobalTracer()
		}
		carrier := opentracing.HTTPHeadersCarrier(c.Request.Header)

		ctx, _ := tr.Extract(opentracing.HTTPHeaders, carrier)

		sp := tr.StartSpan(c.Request.URL.Path, ext.RPCServerOption(ctx))
		ext.HTTPMethod.Set(sp, c.Request.Method)
		ext.HTTPUrl.Set(sp, c.Request.URL.Path)

		ext.Component.Set(sp, "net/http")
		c.Request = c.Request.WithContext(
			opentracing.ContextWithSpan(c.Request.Context(), sp),
		)
		c.Next()
		ext.HTTPStatusCode.Set(sp, uint16(c.Writer.Status()))
		sp.Finish()
	}
}

//GinContext  get context
func GinContext(c *gin.Context) context.Context {
	if c == nil {
		panic("nil context")
	}
	return c.Request.Context()
}

func GinLog(c *gin.Context) *logrus.Entry {
	data := logrus.Fields{
		"x_request_id":      c.GetHeader("x_request_id"),
		"post":              c.Request.PostForm.Encode(),
		"get":               c.Request.URL.RawQuery,
		"url":               c.Request.URL.Path,
		"refer":             c.Request.Referer(),
		"agent":             c.Request.UserAgent(),
		"clientip":          c.ClientIP(),
		"http_header":       headerEncode(c),
		"http_body":         string(ReadBody(c)),
		"http_context_type": c.Request.Header.Get("Content-Type"),
		"type":              0,
	}
	//for _, val := range []string{"sessionid", "userid"} {
	//	item, err := c.Cookie(val)
	//	if err != nil {
	//		continue
	//	}
	//	data[fmt.Sprintf("cookie.%s", val)] = item
	//}
	//var pars struct {
	//	UserId    string `json:"userid" form:"userid"`
	//	SessionId string `json:"sessionid" form:"sessionid"`
	//	Deviceid  string `json:"deviceid" form:"deviceid"`
	//}
	//err := c.Bind(&pars)
	//if err == nil {
	//	data["userId"] = pars.UserId
	//	data["sessionId"] = pars.SessionId
	//}

	return xlog.CtxLog(GinContext(c)).WithFields(data)
}

func headerEncode(c *gin.Context) string {
	if c.Request == nil || c.Request.Header == nil {
		return ""
	}
	return EncodeHTTPHeader(c.Request.Header)
}
func EncodeHTTPHeader(data map[string][]string) string {
	buffer := bytes.NewBufferString("")
	num := 0
	for key, val := range data {
		if num > 0 {
			buffer.WriteString(", ")
		}
		num = num + 1

		buffer.WriteString(key)
		buffer.WriteString("=")
		buffer.WriteString(strings.Join(val, " "))

	}
	return buffer.String()
}
