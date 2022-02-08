package xtools

import (
	"context"
	"fmt"
	"github.com/jinares/xpkg/xlog"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"testing"
)

func TestXErr(t *testing.T) {
	xlog.InitOnlyTracingLog("sss")
	ctx := xlog.StartContext(context.Background(), "ss")
	fmt.Println("trace:", xlog.TraceID(ctx))
	st := XErr(codes.NotFound, "test", true)
	fmt.Println("start:", ErrString(st))
	st = InternalErr(st, "00000000", true)
	fmt.Println("inter:", ErrString(st))
	st = WrapErr(ctx, st)
	fmt.Println("wrap:", ErrString(st))
	ss := FromXErr(st)
	fmt.Println(JSONToStr(ss.Details()))
	for _, it := range ss.Details() {
		switch vv := it.(type) {
		case *epb.RequestInfo:

			fmt.Println("request_id:", vv.RequestId, " :", vv.ServingData)

		}
	}
	st = WrapErr(ctx, st)
	fmt.Println(JSONToStr(FromXErr(st).Details()))
	st = WrapErr(ctx, st)
	fmt.Println(JSONToStr(FromXErr(st).Details()))
}
