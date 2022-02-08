package xerr

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/jinares/xpkg/xtools"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func XErr(code codes.Code, msg string, where ...bool) error {
	//
	//epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	st := status.New(code, msg)
	if len(where) > 0 && where[0] {
		st, err := st.WithDetails(&epb.DebugInfo{
			StackEntries: nil,
			Detail:       xtools.Caller(1, true),
		})
		if err != nil {
			return st.Err()
		}
		return st.Err()
	}

	return st.Err()

}
func FromXErr(err error) *status.Status {
	rpcerr, isok := status.FromError(err)
	if isok {
		return rpcerr
	}
	return rpcerr
}
func WithDetails(err error, info ...proto.Message) error {
	st := FromXErr(err)
	se, err := FromXErr(err).WithDetails(info...)
	if err != nil {
		fmt.Println(fmt.Sprintf("sys-fail:%s", err.Error()))
		return st.Err()
	}
	return se.Err()
}

//String err to string
func ErrString(err error) string {
	if err == nil {
		return ""
	}
	p := FromXErr(err)
	return fmt.Sprintf("%d %s %s", int32(p.Code()), p.Message(), xtools.JSONToStr(p.Details()))
}

//String err to string
func ErrMsg(err error) string {
	if err == nil {
		return ""
	}
	p := FromXErr(err)
	return fmt.Sprintf("%d %s", int32(p.Code()), p.Message())
}
