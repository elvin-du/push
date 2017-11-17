package main

import (
	"fmt"
	"gokit/log"
	"io"
	"net"
	"os"
	"os/signal"
	//	gateMsg "push/gate/message"
	"push/gate/mqtt"
	"push/gate/service/config"
	"push/pb"
	"syscall"

	"google.golang.org/grpc"
)

var (
	defaultServer = &Server{
		UserManager: NewUserManager(),
		Keepalive:   60 * 5, //五分钟检查一次客户端连接情况
	}
)

type Server struct {
	*UserManager
	Keepalive   int64 //单位：秒
	rpcListener net.Listener
	tcpListener net.Listener
}

/*
开始监听RPC端口
*/
func (s *Server) StartRPCServer() {
	var err error = nil
	s.rpcListener, err = net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", config.RPC_SERVICE_PORT))
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		s.rpcListener.Close()
		log.Infof("rpc connection closed")
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		stop := <-ch
		log.Errorf("receive signal '%v'", stop)
		s.Close(stop)
		//		gateMsg.DefaultMessageManager.Sync()
		os.Exit(1)
	}()

	log.Infof("starting rpc server on %d", config.RPC_SERVICE_PORT)

	srv := grpc.NewServer()

	pb.RegisterGateServer(srv, &Gate{})
	srv.Serve(s.rpcListener)
}

//开始监听客户端的连接
func (s *Server) StartTcpServer() {
	var err error = nil
	s.tcpListener, err = net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", config.TCP_PORT))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer func() {
		s.tcpListener.Close()
		log.Infof("tcp connection closed")
	}()

	log.Infof("starting tcp server on %d", config.TCP_PORT)

	for {
		conn, err := s.tcpListener.Accept()
		if nil != err {
			log.Error(err)
			if io.EOF != err {
				return
			}

			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	ses := mqtt.NewSession(conn)
	ses.OnClose(OnClose)
	ses.OnReadPacket(OnRead)
	ses.OnSendPacket(OnSend)
	//启动两个goroutine进行读写
	ses.Start()
}

func (s *Server) Close(reason interface{}) {
	if nil != s.tcpListener {
		s.tcpListener.Close()
	}

	if nil != s.rpcListener {
		s.rpcListener.Close()
	}
}

func (s *Server) PutUser(u *User) {
	s.UserManager.Put(u)
}

func (s *Server) RemoveUser(appID, regID string) *User {
	return s.UserManager.Remove(appID, regID)
}
