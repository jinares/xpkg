package xconn

import (
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

//RetryOption RetryOption
//@times 重试次数
//@retrycode 返回那些错误码则发生重试
func RetryOption(times uint, retrycode ...codes.Code) []grpc.CallOption {
	return []grpc.CallOption{
		grpc_retry.WithMax(times),
		grpc_retry.WithCodes(retrycode...),
	}
}
