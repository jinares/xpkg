package xtools

import (
	"context"
	"errors"
	"fmt"
	"github.com/jinares/xpkg/xlog"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func XErr(code codes.Code, msg string, where ...bool) error {
	//
	//epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	st := status.New(code, msg)
	if len(where) > 0 && where[0] {
		nt, err := st.WithDetails(&epb.DebugInfo{
			StackEntries: nil,
			Detail:       Caller(1, true),
		})
		if err != nil {
			return st.Err()
		}
		return nt.Err()
	}

	return st.Err()

}
func MErr(serr error, code codes.Code, msg string, where ...bool) error {
	if serr == nil {
		serr = errors.New("")
	}
	if msg == "" {
		msg = serr.Error()
	}
	st := status.New(code, msg)
	if len(where) > 0 && where[0] {
		nt, err := st.WithDetails(&epb.DebugInfo{
			StackEntries: nil,
			Detail:       Caller(1, true),
		}, &epb.ErrorInfo{
			Reason:   ErrString(serr),
			Domain:   "",
			Metadata: nil,
		})
		if err != nil {
			return st.Err()
		}
		return nt.Err()
	} else {
		nt, err := st.WithDetails(&epb.ErrorInfo{
			Reason:   ErrString(serr),
			Domain:   "",
			Metadata: nil,
		})
		if err != nil {
			return st.Err()
		}
		return nt.Err()
	}
}

func InternalErr(serr error, msg string, where ...bool) error {
	if serr == nil {
		serr = errors.New("")
	}
	if msg == "" {
		msg = serr.Error()
	}
	st := status.New(codes.Internal, msg)
	if len(where) > 0 && where[0] {
		nt, err := st.WithDetails(&epb.DebugInfo{
			StackEntries: nil,
			Detail:       Caller(1, true),
		}, &epb.ErrorInfo{
			Reason:   ErrString(serr),
			Domain:   "",
			Metadata: nil,
		})
		if err != nil {
			return st.Err()
		}
		return nt.Err()
	} else {
		nt, err := st.WithDetails(&epb.ErrorInfo{
			Reason:   ErrString(serr),
			Domain:   "",
			Metadata: nil,
		})
		if err != nil {
			return st.Err()
		}
		return nt.Err()
	}
}
func WrapErr(ctx context.Context, err error) error {
	if ctx == nil {
		return err
	}
	if err == nil {
		return err
	}
	traceid := xlog.TraceID(ctx)
	if traceid == "" {
		return err
	}
	details := FromXErr(err).Details()
	var hasTrace bool = false
	for _, it := range details {
		switch vv := it.(type) {
		case *epb.RequestInfo:
			if strings.EqualFold(vv.ServingData, "trace") == false {
				continue
			}
			if vv.RequestId == "" {
				continue
			}
			hasTrace = true
			break
		}
	}
	if hasTrace {
		return err
	}
	st, serr := FromXErr(err).WithDetails(&epb.RequestInfo{
		RequestId:   xlog.TraceID(ctx),
		ServingData: "trace",
	})
	if serr != nil {
		return err
	}
	return st.Err()
}
func FromXErr(err error) *status.Status {
	if err == nil {
		return status.New(codes.OK, "ok")
	}
	rpcerr, isok := status.FromError(err)
	if isok {
		return rpcerr
	}
	return rpcerr
}
func Codes(err error) codes.Code {
	return FromXErr(err).Code()
}

//String err to string
func ErrString(err error) string {
	if err == nil {
		return ""
	}
	p := FromXErr(err)
	if len(p.Details()) == 0 {
		return fmt.Sprintf("%d %s", int32(p.Code()), p.Message())
	}
	return fmt.Sprintf("%d %s %s", int32(p.Code()), p.Message(), JSONToStr(p.Details()))
}
func ErrDetail(err error) string {
	if err == nil {
		return ""
	}
	p := FromXErr(err)
	if len(p.Details()) == 0 {
		return ""
	}
	return JSONToStr(p.Details())
}

//String err to string
func ErrMsg(err error) string {
	if err == nil {
		return ""
	}
	p := FromXErr(err)
	return fmt.Sprintf("%d %s", int32(p.Code()), p.Message())
}
