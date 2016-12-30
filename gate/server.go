package main

import (
	"log"
	"net"
	"push/common/server"
	"push/common/util"
	"push/meta"

	"io/ioutil"
	"push/data/client"

	"google.golang.org/grpc"
)

const (
	C_SERVICE_NAME    = "GATE"
	C_SERVICE_VERSION = "1.0"
	C_RPC_PORT        = ":50002"
	C_TCP_PORT        = ":60001"
)

var (
	rpcServer *server.RPCServer
)

/*
开始监听RPC端口
*/
func StartRPCServer() {
	srv := grpc.NewServer()
	meta.RegisterGateServer(srv, &Gate{})

	rpcServer = server.NewRPCServer(util.APP_NAME, C_SERVICE_NAME, C_SERVICE_VERSION, C_RPC_PORT, nil, util.HEARTBEAT_INTERNAL, srv)
	rpcServer.Run()
}

/*
开始监听客户端的连接
*/
func StartTcpServer() {
	l, err := net.Listen("tcp", C_TCP_PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("gate tcp listening at:", C_TCP_PORT)
	//	defer l.Close()

	for {
		conn, err := l.Accept()
		if nil != err {
			log.Println(err)
			continue
		}

		bin, err := ioutil.ReadAll(conn)
		if nil != err {
			log.Println(err)
			continue
		}
		conn.Close()
		log.Println(string(bin))

		resp, err := client.Online(&meta.DataOnlineRequest{UserId: "123", IP: "12.12.151.2"})
		if nil != err {
			log.Println(err)
			continue
		}
		log.Println(resp.String())
	}
}
