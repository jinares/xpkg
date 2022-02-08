package xoidc

import (
	"context"
	"google.golang.org/grpc"
)

func UnaryServerInterceptor(fn OidcAuthFuncHandler) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx, err := Authorize(ctx, info.FullMethod, getMdMap(ctx), fn)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

// GrpcStreamServerInterceptor .
func StreamServerInterceptor(fn OidcAuthFuncHandler) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

		_, err := Authorize(stream.Context(), info.FullMethod, getMdMap(stream.Context()), fn)
		if err != nil {
			return err
		}

		return handler(srv, stream)
	}
}
