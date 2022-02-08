package xconn

import (
	"context"
	"github.com/jinares/xpkg/xlog"
	"github.com/jinares/xpkg/xtools"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"time"
)

//ServerUnaryInvokerLog grpc server call log
func ServerUnaryInvokerLog() grpc.UnaryServerInterceptor {

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		tlog := xlog.CtxLog(ctx)
		tlog = tlog.WithFields(logrus.Fields{
			"par": req,
		})
		var err error
		defer func() {
			ser := xtools.FromXErr(err)
			if ser.Code() == codes.OK {
				tlog.Infof("grpc-server:%s", info.FullMethod)
			} else {
				xlog.Err(err, tlog).Errorf("grpc-server:%s", info.FullMethod)
			}
		}()
		t1 := time.Now()
		resp, err := handler(ctx, req)
		//elapsed := time.Since(t1)
		tlog = xlog.Elapsed(t1, tlog.WithField("resp", resp))
		//tlog = tlog.WithFields(logrus.Fields{
		//	"resp":      resp,
		//	"exec_time": elapsed.Milliseconds(),
		//})
		return resp, err
	}

}
