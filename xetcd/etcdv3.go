package xetcd

import (
	"github.com/jinares/xpkg/xtools"
	"google.golang.org/grpc/codes"
	"strings"
	"time"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/pkg/transport"
)

//EtcdConfig EtcdConfig
type (
	EtcdConfig struct {
		Peers          []string `yaml:"Peers" json:"Peers"`
		TimeoutSeconds int64    `yaml:"TimeoutSeconds" json:"TimeoutSeconds"`
		Root           string   `yaml:"Root" json:"Root"`
		UserName       string   `yaml:"Username" json:"Username"`
		Password       string   `yaml:"Password" json:"Password"`
		Cert           string   `yaml:"Cert" json:"Cert"`
		Key            string   `yaml:"Key" json:"Key"`
		CA             string   `yaml:"CA" json:"CA"`
	}
	EtcdClientv3Handler func() (*clientv3.Client, error)
)

var (
	cliv3 *clientv3.Client
	root  string
)

func GetClientv3() (*clientv3.Client, error) {
	if cliv3 == nil {
		return nil, xtools.XErr(codes.Unimplemented, "没有初始化etcd", true)
	}

	return cliv3, nil
}
func GetRoot() string {
	return strings.TrimSuffix(root, "/")
}

//NewEtcdClientV3 NewEtcdClientV3
func InitEtcdConifg(opt *EtcdConfig) error {
	if cliv3 != nil {
		return xtools.XErr(codes.AlreadyExists, "already")
	}
	if opt.Root == "" || opt.Root == "/" {
		return xtools.XErr(codes.InvalidArgument, "root is empty", true)
	}
	c := clientv3.Config{Username: opt.UserName, Password: opt.Password, Endpoints: opt.Peers, DialTimeout: time.Duration(opt.TimeoutSeconds) * time.Second}
	if opt.Cert != "" {
		tlsInfo := transport.TLSInfo{
			CertFile:      opt.Cert,
			KeyFile:       opt.Key,
			TrustedCAFile: opt.CA,
		}
		tlsConfig, err := tlsInfo.ClientConfig()
		if err != nil {
			panic(err)
		}
		c.TLS = tlsConfig
	}
	client, err := clientv3.New(c)
	if err != nil {
		return err
	}
	cliv3 = client
	root = opt.Root
	return nil

}
