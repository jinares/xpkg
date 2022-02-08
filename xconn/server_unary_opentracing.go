package xconn

import (
	"context"
	"fmt"
	"github.com/jinares/xpkg/xlog"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

//ServerUnaryOpentracing rewrite server's interceptor with open tracing
func ServerUnaryOpentracing(tracer opentracing.Tracer) grpc.UnaryServerInterceptor {

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if tracer == nil {
			return handler(ctx, req)
		}
		//从context中取出metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}
		//从metadata中取出最终数据，并创建出span对象
		spanContext, err := tracer.Extract(opentracing.TextMap, xlog.MDReaderWriter{md})
		if err != nil && err != opentracing.ErrSpanContextNotFound {
			fmt.Errorf("extract from metadata err %v", err)
		}
		//初始化server 端的span
		serverSpan := tracer.StartSpan(
			info.FullMethod,
			ext.RPCServerOption(spanContext),
			opentracing.Tag{Key: string(ext.Component), Value: "gRPC"},
			ext.SpanKindRPCServer,
		)
		defer serverSpan.Finish()
		ctx = opentracing.ContextWithSpan(ctx, serverSpan)
		//将带有追踪的context传入应用代码中进行调用
		return handler(ctx, req)
	}
}
