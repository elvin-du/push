package main

import (
	"log"
	"net"
	"push/common/server"
	"push/common/util"
	"push/gate/mqtt"
	"push/meta"

	"push/data/client"

	"google.golang.org/grpc"
)

const (
	C_SERVICE_NAME    = "GATE"
	C_SERVICE_VERSION = "1.0"
	C_RPC_PORT        = ":50002"

	C_TCP_PORT = ":60001"
)

var (
	defaultServer = &Server{
		Services: make(map[string][]*mqtt.Service),
	}
)

type Server struct {
	Services map[string][]*mqtt.Service //key:userid，一个用户有可能在多台设备上登录
}

/*
开始监听RPC端口
*/
func (s *Server) StartRPCServer() {
	srv := grpc.NewServer()
	meta.RegisterGateServer(srv, &Gate{})

	rpcServer := server.NewRPCServer(util.APP_NAME, C_SERVICE_NAME, C_SERVICE_VERSION, C_RPC_PORT, nil, util.HEARTBEAT_INTERNAL, srv)
	rpcServer.Run()
}

/*
开始监听客户端的连接
*/
func (s *Server) StartTcpServer() {
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

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	svc := mqtt.NewService(conn)

	//获取客户端链接信息
	connMsg, err := svc.GetConnectMessage()
	if nil != err {
		log.Println(err)
		return
	}

	svc.UserId = string(connMsg.ClientId())
	svc.ClientId = string(connMsg.ClientId())

	s.SetService(svc)

	onlineReq := &meta.DataOnlineRequest{}
	onlineReq.ClientId = string(connMsg.ClientId())
	onlineReq.UserId = string(connMsg.ClientId())
	onlineReq.IP = conn.RemoteAddr().String()

	//
	_, err = client.Online(onlineReq)
	if nil != err {
		log.Println(err)
		return
	}

	//TODO 重发
	s.CheckOfflineMsg(onlineReq.UserId)

	//启动两个goroutine进行读写
	err = svc.Start()
	if nil != err {
		log.Println(err)
	}
}

func (s *Server) CheckOfflineMsg(userId string) {
	resp, err := client.GetOfflineMsgs(userId)
	if nil != err {
		log.Println(err)
	}

	//TODO
	svcs := s.Services[userId]
	for _, v := range svcs {
		for _, v2 := range resp.Items {
			go v.Push([]byte(v2.Content))
		}
	}
}

func (s *Server) SetService(svc *mqtt.Service) {
	s.Services[svc.UserId] = append(s.Services[svc.UserId], svc)
}
