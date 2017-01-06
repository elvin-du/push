package main

import (
	"push/common/server"
	"push/common/util"

	"push/meta"

	"google.golang.org/grpc"
)

const (
	C_SERVICE_NAME    = "SESSION"
	C_SERVICE_VERSION = "1.0"
	C_RPC_PORT        = ":50003"
)

func StartRPCServer() {
	srv := grpc.NewServer()
	meta.RegisterSessionServer(srv, &Session{})

	rpcServer := server.NewRPCServer(util.APP_NAME, C_SERVICE_NAME, C_SERVICE_VERSION, C_RPC_PORT, nil, util.HEARTBEAT_INTERNAL, srv)
	rpcServer.Run()
}
