package xconn

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"time"
)

var (
	kasp = keepalive.ServerParameters{
		// If a client is idle for 30 seconds, send a GOAWAY,
		MaxConnectionIdle: 30 * time.Second,
		//If any connection is alive for more than 60 seconds, send a GOAWAY
		MaxConnectionAge: 1 * time.Minute,
		// Allow 5 seconds for pending RPCs to complete before forcibly closing connections
		MaxConnectionAgeGrace: 5 * time.Second,
		//Ping the client if it is idle for 10 seconds to ensure the connection is still active
		Time: 10 * time.Second,
		// Wait 1 second for the ping ack before assuming the connection is dead
		Timeout: 1 * time.Second,
	}
	kaep = keepalive.EnforcementPolicy{
		MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
		PermitWithoutStream: true,            // Allow pings even when there are no active streams
	}
)

//SeverKeepAliveEnforcementPolicy EnforcementPolicy
func SeverKeepAliveEnforcementPolicy() grpc.ServerOption {
	return grpc.KeepaliveEnforcementPolicy(kaep)
}

//ServerKeepAlive KeepAliveServer
func ServerKeepAlive() grpc.ServerOption {
	return grpc.KeepaliveParams(kasp)
}
