package xconn

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"

	//"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	//"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/go-grpc-middleware/validator"
)

func AdditionalUnaryServerOptWithBase(item ...grpc.UnaryServerInterceptor) []grpc.UnaryServerInterceptor {
	ret := []grpc.UnaryServerInterceptor{
		grpc_validator.UnaryServerInterceptor(),
		grpc_ctxtags.UnaryServerInterceptor(),
	}
	return append(ret, item...)
}

//AdditionalStreamServerOptWithBase add base
func AdditionalStreamServerOptWithBase(item ...grpc.StreamServerInterceptor) []grpc.StreamServerInterceptor {
	ret := []grpc.StreamServerInterceptor{
		grpc_validator.StreamServerInterceptor(),
		grpc_ctxtags.StreamServerInterceptor(),
	}
	return append(ret, item...)
}
func CreateWithUnaryServerChain(item ...grpc.UnaryServerInterceptor) grpc.ServerOption {
	return grpc_middleware.WithUnaryServerChain(item...)
}

//CreateWithStreamServerChain  WithStreamServerChain
func CreateWithStreamServerChain(item ...grpc.StreamServerInterceptor) grpc.ServerOption {
	return grpc_middleware.WithStreamServerChain(item...)
}
