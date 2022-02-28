module github.com/jinares/xpkg

go 1.14

//go.etcd.io/bbolt => github.com/coreos/bbolt v1.3.5
//replace go.etcd.io/etcd => github.com/etcd-io/etcd v3.3.22+incompatible

require (
	github.com/allegro/bigcache/v2 v2.2.5 // indirect
	github.com/allegro/bigcache/v3 v3.0.0
	github.com/asaskevich/EventBus v0.0.0-20200907212545-49d423059eef // indirect
	github.com/astaxie/beego v1.12.3
	github.com/beltran/gohive v1.4.0
	github.com/bluele/gcache v0.0.0-20190518031135-bc40bd653833
	github.com/codahale/hdrhistogram v0.0.0-20161010025455-3a0bb77429bd // indirect
	github.com/coreos/go-oidc/v3 v3.0.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/doug-martin/goqu/v9 v9.18.0
	github.com/ethereum/go-ethereum v1.10.12
	github.com/gin-gonic/gin v1.7.2
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8
	github.com/go-mysql-org/go-mysql v1.3.0
	github.com/go-oauth2/oauth2/v4 v4.3.0
	github.com/golang/protobuf v1.5.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.1-0.20190118093823-f849b5445de4
	github.com/grpc-ecosystem/grpc-gateway v1.11.2 // indirect
	github.com/myzhan/boomer v1.6.0
	github.com/opentracing/opentracing-go v1.1.0
	github.com/sirupsen/logrus v1.7.0
	github.com/uber/jaeger-client-go v2.20.1+incompatible
	github.com/uber/jaeger-lib v2.2.0+incompatible // indirect
	github.com/zeromq/goczmq v4.1.0+incompatible // indirect
	github.com/zeromq/gomq v0.0.0-20201031135124-cef4e507bb8e // indirect
	github.com/zeromq/gomq/zmtp v0.0.0-20201031135124-cef4e507bb8e // indirect
	go.etcd.io/etcd v0.0.0-20191011172313-6d8052314b9e
	golang.org/x/net v0.0.0-20210805182204-aaa1db679c0d
	golang.org/x/text v0.3.6
	google.golang.org/api v0.52.0
	google.golang.org/genproto v0.0.0-20210722135532-667f2b7c528f
	google.golang.org/grpc v1.39.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/yaml.v2 v2.4.0
)

replace (
	go.etcd.io/etcd => github.com/etcd-io/etcd v0.0.0-20191021022006-5dc12f27251a
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)
