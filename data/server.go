package main

import (
	"push/common/server"
	"push/common/util"
	"push/gate/service/config"
	"push/meta"

	"google.golang.org/grpc"
)

func StartRPCServer() {
	srv := grpc.NewServer()
	meta.RegisterDataServer(srv, &Data{})

	server.NewRPCServer(
		util.APP_NAME,
		config.RpcServiceName,
		config.RpcServiceVersion,
		config.RpcServicePort,
		nil,
		util.HEARTBEAT_INTERNAL,
		srv,
	).Run()
}
