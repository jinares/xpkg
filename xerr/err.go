package xerr

import (
	"fmt"
	"github.com/jinares/xpkg/xtools"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	XError struct {
		Code    codes.Code `json:"code"`
		Message string     `json:"message"`
		Where   string     `json:"where"`
	}
)

func (se *XError) GRPCStatus() *status.Status {
	return status.New(codes.Code(se.Code), se.Message)
}

/*
Error() string
*/
func (p *XError) Error() string {
	return fmt.Sprintf("%d %s", (p.Code), p.Message)
}

/*
Error() string
*/
func (p *XError) String() string {
	return fmt.Sprintf("%s %s", p.Error(), p.Where)
}
func String(err error) string {
	if err == nil {
		return ""
	}
	if val, isok := (err).(*XError); isok {
		return val.String()
	}
	return err.Error()
}

func XErr(code codes.Code, msg string, where ...bool) *XError {

	if len(where) > 0 && where[0] {
		return &XError{
			Code:    code,
			Message: msg,
			Where:   xtools.Caller(1, true), // caller(1, true),
		}
	}
	return &XError{
		Code:    code,
		Message: msg,
	}

}
func FromXErr(err error) *XError {
	if err == nil {
		return XErr(codes.OK, "ok")
	}
	if xe, isok := err.(*XError); isok {
		return xe
	}
	rpcerr, isok := status.FromError(err)
	if isok {
		return &XError{Code: rpcerr.Code(), Message: rpcerr.Message()}
	}
	msg := err.Error()
	index := strings.Index(msg, " ")
	if index < 1 {
		return &XError{
			Code:    codes.Unknown,
			Message: msg,
		}
	}
	codestr := msg[0:index]

	icode, isok := xtools.IntVal(codestr)
	if isok == false {
		return &XError{
			Code:    codes.Unknown,
			Message: msg,
		}
	}
	return &XError{Code: codes.Code(icode), Message: msg[index:]}
}
