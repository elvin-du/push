package main

import (
	"push/common/server"
	"push/common/util"
	"push/data/service/config"
	"push/meta"

	"google.golang.org/grpc"
)

const (
	RPC_SERVICE_NAME    = "DATA"
	RPC_SERVICE_VERSION = "1.0.0"
)

func StartRPCServer() {
	srv := grpc.NewServer()
	meta.RegisterDataServer(srv, &Data{})

	server.NewRPCServer(
		util.APP_NAME,
		RPC_SERVICE_NAME,
		RPC_SERVICE_VERSION,
		config.SERVER_IP,
		config.RPC_SERVICE_PORT,
		nil,
		util.HEARTBEAT_INTERNAL,
		srv,
	).Run()
}
