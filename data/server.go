package main

import (
	"log"

	"net"
	"push/common/etcd"
	"push/common/server"
	"push/common/util"

	"push/meta"

	"google.golang.org/grpc"
)

const (
	C_SERVICE_NAME    = "DATA"
	C_SERVICE_VERSION = "1.0"
	C_RPC_PORT            = ":50001"
)

func StartRPCServer() {
	l, err := net.Listen("tcp", C_RPC_PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("gate rpc listening at:", C_RPC_PORT)
	defer l.Close()

	srv := grpc.NewServer()
	meta.RegisterDataServer(srv, &Data{})

	cli, err := etcd.NewETCDClient(util.EtcdEndpoints, util.DialTimeout, util.RequestTimeout)
	if nil != err {
		log.Fatalln(err)
	}

	err = server.NewServer(util.APPName, C_SERVICE_NAME, C_SERVICE_VERSION, C_RPC_PORT, nil, cli, util.HeartbeatInterval).Start()
	if nil != err {
		log.Fatalln(err)
	}

	err = srv.Serve(l)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
}
