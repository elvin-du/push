package main

import (
	"errors"
	"hscore/log"
	"net"
	"push/common/server"
	"push/common/util"
	dataCli "push/data/client"
	"push/gate/mqtt"
	"push/gate/service/config"
	"push/meta"
	"strings"
	"time"

	"github.com/surgemq/message"
	"google.golang.org/grpc"
)

var (
	defaultServer = &Server{
		Services: make(map[string]*mqtt.Service),
	}
)

type Server struct {
	Services map[string]*mqtt.Service //key:userid，一个用户不可以同时在多台设备上登录
}

/*
开始监听RPC端口
*/
func (s *Server) StartRPCServer() {
	srv := grpc.NewServer()
	meta.RegisterGateServer(srv, &Gate{})

	rpcServer := server.NewRPCServer(
		util.APP_NAME,
		config.RpcServiceName,
		config.RpcServiceVersion,
		config.RpcServicePort,
		nil,
		util.HEARTBEAT_INTERNAL,
		srv)
	rpcServer.Run()
}

/*
开始监听客户端的连接
*/
func (s *Server) StartTcpServer() {
	l, err := net.Listen("tcp", config.TcpPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Debugln("gate tcp listening at:", config.TcpPort)
	defer l.Close()

	for {
		conn, err := l.Accept()
		if nil != err {
			log.Error(err)
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	svc := mqtt.NewService(conn)
	svc.Keepalive = time.Minute * 5

	defer func() {
		s.DelService(svc)
	}()

	//获取客户端链接信息
	connMsg, err := svc.GetConnectMessage()
	if nil != err {
		log.Error(err)
		return
	}
	log.Infof("clientid:%s connected", connMsg.ClientId())

	//启动两个goroutine进行读写
	go svc.Run()

	svc.UserId = string(connMsg.Username())
	svc.ClientId = string(connMsg.ClientId())
	err = svc.SetWriteDeadline(svc.Keepalive)
	if nil != err {
		log.Error(err)
		return
	}

	connAckMsg := message.NewConnackMessage()

	//合法性检验
	err = s.Auth(svc.UserId, string(connMsg.Password()))
	if nil != err {
		connAckMsg.SetReturnCode(message.ErrNotAuthorized)
		err = svc.Write(connAckMsg)
		if nil != err {
			log.Error(err)
			return
		}

		return
	}

	platform, err := s.ParseClientId(svc.ClientId)
	if nil != err {
		connAckMsg.SetReturnCode(message.ErrNotAuthorized)
		err = svc.Write(connAckMsg)
		if nil != err {
			log.Error(err)
			return
		}

		return
	}

	//回应connect消息
	connAckMsg.SetReturnCode(message.ConnectionAccepted)
	err = svc.Write(connAckMsg)
	if nil != err {
		log.Error(err)
		return
	}

	//每一个新用户链接抽象为一个service，并把其保存到相应的server实例中
	s.SetService(svc)

	s.TryKickOff(svc.ClientId, svc.UserId)
	s.CheckOfflineMsg(svc.UserId)

	gateIp, err := util.LocalIP()
	if nil != err {
		log.Error(err)
		return
	}

	onlineReq := &meta.DataOnlineRequest{}
	onlineReq.ClientId = svc.ClientId
	onlineReq.UserId = svc.UserId
	onlineReq.GateIp = gateIp
	onlineReq.GatePort = config.RpcServicePort
	onlineReq.Platform = platform

	//
	_, err = dataCli.Online(onlineReq)
	if nil != err {
		log.Error(err)
		return
	}

	log.Infoln("should not show")
}

func (s *Server) CheckOfflineMsg(userId string) {
	resp, err := dataCli.GetOfflineMsgs(userId)
	if nil != err {
		log.Error(err)
	}

	log.Debugf("found %d offline msg for %s", len(resp.Items), userId)
	svc := s.Services[userId]
	for _, v2 := range resp.Items {
		go svc.Push(uint16(v2.PacketId), []byte(v2.Content))
	}
}

/*
如果此用户已经在别的设备上登录，断开其链接
*/
func (s *Server) TryKickOff(clientId, userId string) {
	//TODO
}

func (s *Server) SetService(svc *mqtt.Service) {
	s.Services[svc.UserId] = svc
}

func (s *Server) DelService(svc *mqtt.Service) {
	if nil != svc && nil != s.Services {
		delete(s.Services, svc.UserId)
	}
}

/*
因为需要辨别每次的链接是否是同一个手机，所以需要手机根据硬件生成一个唯一标识来当作clientId,
clientId格式：OS系统-手机硬件唯一对应标识．例如：android-123abb
*/
func (s *Server) ParseClientId(clientId string) (string, error) {
	platform := strings.TrimLeft(clientId, "-")
	if "android" == platform || "ios" == platform {
		return platform, nil
	}

	log.Errorf("clientId:%s", clientId)
	return "", errors.New("clientId invalid")
}

//TODO
func (s *Server) Auth(userId, token string) error {
	return nil
}
