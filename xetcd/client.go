package xetcd

import (
	"context"
	"fmt"
	"github.com/jinares/xpkg/xtools"
	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc/codes"
)

func LinkKey(key string) string {
	return fmt.Sprintf("%s/%s", GetRoot(), key)
}
func Get(ctx context.Context, key string) (string, error) {
	cli, err := GetClientv3()
	if err != nil {
		return "", err
	}

	res, err := cli.Get(ctx, LinkKey(key))
	if err != nil {
		return "", err
	}
	if res.Count < 1 {
		return "", xtools.XErr(codes.NotFound, "empty")
	}
	return string(res.Kvs[0].Value), nil

}

func GetWithFirstCreate(ctx context.Context, key string) (string, error) {
	cli, err := GetClientv3()
	if err != nil {
		return "", err
	}

	res, err := cli.Get(ctx, LinkKey(key), clientv3.WithLastCreate()...)
	if err != nil {
		return "", err
	}
	if res.Count < 1 {
		return "", xtools.XErr(codes.NotFound, "empty")
	}
	return string(res.Kvs[0].Value), nil

}

func Set(ctx context.Context, key, val string) error {
	cli, err := GetClientv3()
	if err != nil {
		return err
	}
	_, err = cli.Put(ctx, LinkKey(key), val)
	return err
}
