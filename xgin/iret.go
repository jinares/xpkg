package xgin

import (
	"github.com/gin-gonic/gin"
	"github.com/jinares/xpkg/xlog"
	"github.com/jinares/xpkg/xtools"
	"google.golang.org/grpc/codes"
)

type (
	IRet interface {
		GetCode() codes.Code
		GetRet() (string, ContentType)
		SetTrace(trace string) IRet
	}

	Ret struct {
		Ret   codes.Code  `json:"ret"`
		Msg   string      `json:"msg,omitempty"`
		Data  interface{} `json:"data,omitempty"`
		Trace string      `json:"trace,omitempty"`
	}
)

func Error(c *gin.Context, err error) IRet {
	trace := xlog.TraceID(GinContext(c))
	ret := NewRet(err, nil)
	ret.SetTrace(trace)
	return ret
}
func TraceID(c *gin.Context) string {
	return xlog.TraceID(GinContext(c))
}
func Trace(c *gin.Context, ret IRet) IRet {
	if ret == nil {
		return ret
	}
	return ret.SetTrace(xlog.TraceID(GinContext(c)))
}
func Succ(c *gin.Context, data interface{}) IRet {
	return NewRet(nil, data)
}
func NewRet(err error, data interface{}) IRet {
	ret := &Ret{Data: data}
	if err != nil {
		xe := xtools.FromXErr(err)
		ret.Msg = xe.Message()
		ret.Ret = xe.Code()
		if xe.Code() != codes.OK {
			ret.Data = nil
			//ret.Tips=""
		}
	}
	return ret
}
func (c *Ret) GetCode() codes.Code {
	return c.Ret
}
func (c *Ret) GetRet() (string, ContentType) {
	return xtools.JSONToStr(c), APPLICATION_JSON
}

func (c *Ret) SetTrace(str string) IRet {
	c.Trace = str
	return c
}

/*

func (h *Msg) MarshalJSON() ([]byte, error) {

	return []byte("\"" + h.Str() + "\""), nil
}
func (h *Msg) UnmarshalJSON(data []byte) error {
	str := string(data)
	h.data = strings.Split(str, "-")
	return nil
}
*/
