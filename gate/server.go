package main

import (
	"log"
	"net"
	"push/common/server"
	"push/common/util"
	"push/gate/service"
	"push/meta"

	"github.com/surgemq/message"
	"google.golang.org/grpc"

	"io"
)

const (
	C_SERVICE_NAME    = "GATE"
	C_SERVICE_VERSION = "1.0"
	C_RPC_PORT        = ":50002"

	C_TCP_PORT = ":60001"
)

/*
开始监听RPC端口
*/
func StartRPCServer() {
	srv := grpc.NewServer()
	meta.RegisterGateServer(srv, &Gate{})

	rpcServer := server.NewRPCServer(util.APP_NAME, C_SERVICE_NAME, C_SERVICE_VERSION, C_RPC_PORT, nil, util.HEARTBEAT_INTERNAL, srv)
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
	defer l.Close()

	for {
		conn, err := l.Accept()
		if nil != err {
			log.Println(err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	resp := message.NewConnackMessage()
}

func getConnectMsg(r io.Reader) (*message.ConnectMessage, error) {
	var (
		header = make([]byte, 1)
	)

	n, err := r.Read(header)
	if nil != err {
		log.Println(err)
		return nil, err
	}

}
