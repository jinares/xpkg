package xconn

import (
	"context"
	"github.com/jinares/xpkg/xlog"
	"github.com/jinares/xpkg/xtools"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"runtime/debug"
)

//ServerUnaryRecoveryInterceptor RecoveryInterceptor,
func ServerUnaryRecoveryInterceptor(
	ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (ret interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			//debug.PrintStack()
			stack := debug.Stack()
			xlog.CtxLog(ctx).WithFields(logrus.Fields{
				"stack":      string(stack),
				"err_string": xtools.JSONToStr(e),
			}).Errorf("gprc-recover:%s", info.FullMethod)
			err = xtools.XErr(codes.Internal, "please try again later")
		}
	}()

	return handler(ctx, req)
}
