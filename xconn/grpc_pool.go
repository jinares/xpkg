package xconn

import (
	"errors"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/jinares/xpkg/xtools"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
	"math"
	"net/url"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrConnShutdown = errors.New("grpc-conn-shutdown")

	defaultClientPoolCap = 5
	pool                 sync.RWMutex
	connpool             = map[string]*ClientPool{}

	kacp = keepalive.ClientParameters{
		Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
		Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
		PermitWithoutStream: true,             // send pings even without active streams
	}
)

type (
	ClientPool struct {
		option   []grpc.DialOption
		capacity int64
		next     int64
		target   string

		sync.Mutex

		conns []*grpc.ClientConn
	}
)

func RegisterPool(target string, cli *ClientPool) error {
	if target == "" || cli == nil {
		return status.Error(codes.InvalidArgument, "client pool invalid argument")
	}
	pool.Lock()
	defer pool.Unlock()
	connpool[target] = cli
	return nil
}
func CloseGrpcConn() {
	pool.Lock()
	defer pool.Unlock()
	for _, item := range connpool {
		item.Close()
	}
	connpool = map[string]*ClientPool{}
}
func GrpcConnPoolv2(target string, cap int, opt ...grpc.DialOption) (*grpc.ClientConn, error) {
	if target == "" {
		return nil, xtools.XErr(codes.Internal, "", true)
	}
	pool.RLock()
	v, isok := connpool[target]
	pool.RUnlock()
	if isok {
		return v.GetConn()
	}
	if len(opt) < 1 {

		serviceName := os.Getenv("JAEGER_SERVICE_NAME")
		opt = []grpc.DialOption{
			grpc.WithInsecure(),
			grpc.WithUnaryInterceptor(
				grpc_middleware.ChainUnaryClient(
					ClientUnaryOpentracingAndTimeOut(serviceName, 6*time.Second),
					ClientUnaryOpentracing(),
				),
			),
			grpc.WithKeepaliveParams(kacp),
		}
	}
	pool.Lock()
	defer pool.Unlock()
	cli, isok := connpool[target]
	if isok {
		return cli.GetConn()
	}
	cli = NewClientPool(target, cap, opt...)
	connpool[target] = cli

	return cli.GetConn()
}
func GrpcConnPool(target string, opt ...grpc.DialOption) (*grpc.ClientConn, error) {
	return GrpcConnPoolv2(target, defaultClientPoolCap, opt...)

}
func NewClientPool(target string, poolsize int, option ...grpc.DialOption) *ClientPool {
	if poolsize < 1 {
		poolsize = defaultClientPoolCap
	}
	u, err := url.Parse(target)
	if err == nil && u.Host != "" {
		target = u.Host
	}

	return &ClientPool{
		target:   target,
		conns:    make([]*grpc.ClientConn, poolsize),
		capacity: int64(poolsize),
		option:   option,
	}
}

func (cc *ClientPool) checkState(conn *grpc.ClientConn) error {
	if conn == nil {
		return ErrConnShutdown
	}
	state := conn.GetState()
	switch state {
	case connectivity.TransientFailure, connectivity.Shutdown:
		return ErrConnShutdown
	}
	return nil
}

func (cc *ClientPool) GetConn() (*grpc.ClientConn, error) {
	var (
		idx  int64
		next int64
		err  error
	)

	next = atomic.AddInt64(&cc.next, 1)
	if next >= math.MaxInt64 {
		cc.next = 0
	}
	idx = next % cc.capacity
	conn := cc.conns[idx]
	if cc.checkState(conn) == nil {
		return conn, nil
	}

	cc.Lock()
	defer cc.Unlock()
	conn = cc.conns[idx]
	if cc.checkState(conn) == nil {
		return conn, nil
	}
	if conn != nil {
		conn.Close()
	}
	conn, err = cc.connect()
	if err != nil {
		return nil, err
	}

	cc.conns[idx] = conn
	return conn, nil
}
func (cc *ClientPool) connect() (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(cc.target, cc.option...)
	if err != nil {
		return nil, xtools.XErr(codes.Internal, err.Error(), true)
	}
	return conn, nil
}

func (cc *ClientPool) Close() {
	cc.Lock()
	defer cc.Unlock()

	for _, conn := range cc.conns {
		if conn == nil {
			continue
		}

		conn.Close()
	}
	cc.conns = make([]*grpc.ClientConn, cc.capacity)
}
