package xetcd

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.etcd.io/etcd/clientv3"
)

type (
	Option struct {
		Key     string
		Convert ConvertFunc
	}
	tmpOption struct {
		Op   Option
		Data map[string]string
	}

	ConvertFunc      func(data interface{}) error
	WatchNodeFunc    func(val string) error
	WatchDirNodeFunc func(key, val string) error
)

func WatchNode(ctx context.Context, path string, action WatchNodeFunc) error {
	client, err := GetClientv3()
	if err != nil {
		return err
	}
	val, err := getNode(ctx, client, path)
	if err == nil && val != "" {
		action(val)
	}
	go func() {
		for {
			cc := client.Watch(ctx, path)
			for wres := range cc {
				fmt.Println(wres)
				val, err := getNode(ctx, client, path)
				if err == nil && val != "" {
					action(val)
				}
			}
		}
	}()
	return nil

}

func WatchDirNode(ctx context.Context, path string, action WatchDirNodeFunc) error {
	client, err := GetClientv3()
	if err != nil {
		return err
	}
	val, err := getDirNode(ctx, client, path)
	if err == nil {
		for k, v := range val {
			action(k, v)
		}
	}
	go func() {
		for {
			cc := client.Watch(ctx, path, clientv3.WithPrefix())
			for wresp := range cc {

				for _, ev := range wresp.Events {

					if len(ev.Kv.Value) < 1 {
						continue
					}
					k, sk := split(path, string(ev.Kv.Key))
					if sk != "" {
						fmt.Println(ev.Kv)
						continue
					}
					action(k, string(ev.Kv.Value))

				}
			}
		}
	}()
	return nil

}

//GetNode
func getNode(ctx context.Context, client *clientv3.Client, key string) (string, error) {
	res, err := client.Get(ctx, key)
	if err != nil {
		return "", err
	}
	if res.Count < 1 {
		return "", errors.New("empty")
	}
	return string(res.Kvs[0].Value), nil
}
func getDirNode(ctx context.Context, client *clientv3.Client, path string) (map[string]string, error) {

	res, err := client.Get(ctx, path, clientv3.WithPrefix())
	if err != nil {
		return map[string]string{}, err
	} else {
		data := map[string]string{}
		for _, v := range res.Kvs {
			key, subkey := split(path, string(v.Key))
			if key == "" {
				continue
			}
			if subkey != "" {
				continue
			}
			if strings.HasSuffix(string(v.Key), "/") {
				continue
			}
			data[key] = string(v.Value)
		}
		return data, nil
	}
	return nil, errors.New("")
}

func match(root string, mop map[string]Option, res *clientv3.GetResponse) {
	tmp := map[string]tmpOption{}
	for _, v := range res.Kvs {
		key, subkey := split(root, string(v.Key))
		op, isok := mop[key]
		if isok == false {
			continue
		}
		if op.Convert == nil {
			continue
		}
		if subkey == "" {
			op.Convert(string(v.Value))
		} else {
			if vop, isok := tmp[key]; isok {
				vop.Data[subkey] = string(v.Value)
			} else {
				tmp[key] = tmpOption{Op: op, Data: map[string]string{
					subkey: string(v.Value),
				}}
			}
		}
	}
	for _, val := range tmp {
		val.Op.Convert(val.Data)
	}
}
func split(root, path string) (key string, subkey string) {

	data := strings.TrimPrefix(path, root)
	data = strings.TrimSuffix(strings.TrimPrefix(data, "/"), "/")
	if data == "" {
		return
	}
	arr := strings.Split(data, "/")
	if len(arr) < 1 || len(arr) > 2 {
		return
	}
	key = arr[0]
	subkey = strings.Join(arr[1:], "/")
	return
}
