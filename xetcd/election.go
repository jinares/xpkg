package xetcd

import (
	"context"
	"github.com/jinares/xpkg/xlog"
	"go.etcd.io/etcd/clientv3/concurrency"
)

//NewElection NewElection
func Election(ctx context.Context, key, val string, fn MutexHandler, sopts ...concurrency.SessionOption) error {
	cli, err := GetClientv3()
	if err != nil {
		return err
	}
	session, err := concurrency.NewSession(cli, sopts...)
	if err != nil {
		return err
	}
	defer session.Close()
	e1 := concurrency.NewElection(session, LinkKey(key))
	if err := e1.Campaign(context.Background(), val); err != nil {
		return err
	}
	defer func() {
		err := e1.Resign(ctx)
		if err != nil {
			xlog.Err(err).Error("etcd-election-fail")
		}
	}()
	return fn()
}
