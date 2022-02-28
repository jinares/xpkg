package xetcd

import (
	"context"
	"fmt"
	"github.com/jinares/xpkg/xlog"
	"github.com/jinares/xpkg/xtools"
	"google.golang.org/grpc/codes"
	"runtime/debug"

	"go.etcd.io/etcd/clientv3/concurrency"
)

type (
	MutexHandler func() error
)

//Mutex Mutex
func Mutex(key string, f MutexHandler, sopts ...concurrency.SessionOption) error {
	cli, err := GetClientv3()
	if err != nil {
		return err
	}
	session, err := concurrency.NewSession(cli, sopts...)
	if err != nil {
		return err
	}

	defer session.Close()
	fmt.Println("linkkey:", LinkKey(key))
	mu := concurrency.NewMutex(session, LinkKey(key))

	if err := mu.Lock(context.TODO()); err != nil {
		return err
	}
	defer func() {
		err := mu.Unlock(context.TODO())
		if err != nil {
			xlog.Err(err).Error("etcd-mutex-unlock-fail")

		}
	}()

	return f()

}
func MutexLoopHandler(key string, fn MutexHandler, sopts ...concurrency.SessionOption) error {
	defer func() {
		if e := recover(); e != nil {
			//debug.PrintStack()
			stack := debug.Stack()

			xlog.Infof("mutex-recover-fail:msg:  %v stack: %s", e, string(stack))
		}
	}()
	for {
		err := Mutex(key, fn, sopts...)
		if err != nil {
			xlog.Err(err).Error("mutex-fail")
		}
	}
	return xtools.XErr(codes.Internal, "mutex-fail")
}
