module github.com/jinares/xpkg

go 1.14

//go.etcd.io/bbolt => github.com/coreos/bbolt v1.3.5
//replace go.etcd.io/etcd => github.com/etcd-io/etcd v3.3.22+incompatible

require (
	github.com/bluele/gcache v0.0.0-20190518031135-bc40bd653833
	github.com/codahale/hdrhistogram v0.0.0-20161010025455-3a0bb77429bd // indirect
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8
	github.com/golang/protobuf v1.3.2
	github.com/grpc-ecosystem/grpc-gateway v1.11.2 // indirect
	github.com/opentracing/opentracing-go v1.1.0
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.4.0 // indirect
	github.com/uber/jaeger-client-go v2.20.1+incompatible
	github.com/uber/jaeger-lib v2.2.0+incompatible // indirect
	go.etcd.io/etcd v0.0.0-20191011172313-6d8052314b9e
	golang.org/x/text v0.3.2
	golang.org/x/time v0.0.0-20181108054448-85acf8d2951c // indirect
	google.golang.org/genproto v0.0.0-20191108220845-16a3f7862a1a // indirect
	google.golang.org/grpc v1.24.0
	gopkg.in/yaml.v2 v2.2.4
)

replace go.etcd.io/etcd => github.com/etcd-io/etcd v0.0.0-20191021022006-5dc12f27251a
