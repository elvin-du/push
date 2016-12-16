package main

import (
	"log"

	"net"

	"push/common/etcd"
	"push/common/server"
	"push/common/util"

	"golang.org/x/net/context"

	"push/meta"

	"io/ioutil"

	"google.golang.org/grpc"
)

const (
	C_SERVICE_NAME    = "GATE"
	C_SERVICE_VERSION = "1.0"
	C_RPC_PORT        = ":50002"
	C_TCP_PORT        = ":60001"
)

var (
	rpcServer *server.Server
)

func init() {
	log.SetFlags(log.Lshortfile)
}

/*
开始监听RPC端口
*/
func StartRPCServer() {
	l, err := net.Listen("tcp", C_RPC_PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("gate rpc listening at:", C_RPC_PORT)
	//	defer l.Close()

	srv := grpc.NewServer()
	meta.RegisterGateServer(srv, &Gate{})

	cli, err := etcd.NewETCDClient(util.EtcdEndpoints, util.DialTimeout, util.RequestTimeout)
	if nil != err {
		log.Fatalln(err)
	}

	rpcServer = server.NewServer(util.APPName, C_SERVICE_NAME, C_SERVICE_VERSION, C_RPC_PORT, nil, cli, util.HeartbeatInterval)
	err = rpcServer.Start()
	if nil != err {
		log.Fatalln(err)
	}

	err = srv.Serve(l)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
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

		//服务发现
		ip, port, err := rpcServer.Client.Get("YD", "DATA", "1.0")
		if nil != err {
			log.Println(err)
			continue
		}
		log.Println(ip, port)

		target := ip + port
		log.Print(target)
		cliConn, err := grpc.Dial(target, grpc.WithInsecure())
		if nil != err {
			log.Println(err)
			continue
		}
		//        defer cliConn.Close()
		client := meta.NewDataClient(cliConn)

		resp, err := client.Online(context.Background(), &meta.DataOnlineRequest{UserId: "123", IP: string(bin)})
		if nil != err {
			log.Println(err)
			continue
		}
		log.Println(resp.String())

		log.Println(string(bin))
	}
}
