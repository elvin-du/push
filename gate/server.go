package main

import (
	//	"errors"
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
		Services:  make(map[string]*mqtt.Service),
		Keepalive: 60 * 5, //五分钟检查一次客户端连接情况
	}
)

type Server struct {
	Services  map[string]*mqtt.Service //key:userid，一个用户不可以同时在多台设备上登录
	Keepalive int64                    //单位：秒
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

//开始监听客户端的连接
func (s *Server) StartTcpServer() {
	l, err := net.Listen("tcp", config.TcpPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer l.Close()

	log.Infoln("gate tcp listening at:", config.TcpPort)

	//开始监测所有连接本server的客户端的连接状况
	go s.CronEvery()

	for {
		conn, err := l.Accept()
		if nil != err {
			log.Error(err)
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) checkConnection(svc *mqtt.Service) (err error) {
	var connMsg *message.ConnectMessage
	connAckMsg := message.NewConnackMessage()
	defer func() {
		//回应connect消息
		if nil != err {
			connAckMsg.SetReturnCode(message.ErrNotAuthorized)
			err2 := svc.SendMsg(connAckMsg)
			if nil != err2 {
				log.Error(err2)
			}

			//直接关闭连接
			if nil != svc.Conn {
				err2 = svc.Conn.Close()
				if nil != err2 {
					log.Error(err2)
				}
			}
		}
	}()

	//获取客户端链接信息
	connMsg, err = svc.GetConnectMessage()
	if nil != err {
		log.Error(err)
		return err
	}
	log.Infoln("come to connect,clientid:", string(connMsg.ClientId()))

	svc.UserId = string(connMsg.Username())
	svc.ClientId = string(connMsg.ClientId())

	//合法性检验
	err = s.Auth(svc.UserId, string(connMsg.Password()))
	if nil != err {
		log.Error(err)
		return err
	}

	var platform string
	platform, err = s.ParseClientId(svc.ClientId)
	if nil != err {
		log.Error(err)
		return err
	}
	svc.Platform = platform

	//连接成功
	connAckMsg.SetReturnCode(message.ConnectionAccepted)
	err = svc.SendMsg(connAckMsg)
	if nil != err {
		log.Error(err)
		return err
	}

	svc.SetTouchTime(time.Now().Unix())
	//启动两个goroutine进行读写
	svc.Run()

	return nil
}

func (s *Server) handleConnection(conn net.Conn) {
	svc := mqtt.NewService(conn)
	svc.SetWriteTimeout(time.Second * 30)
	svc.SetReadTimeout(time.Minute * 10)

	err := s.checkConnection(svc)
	if nil != err {
		log.Errorln(err)
		return
	}

	//每一个新用户链接抽象为一个service，并把其保存到相应的server实例中
	s.SetService(svc)

	//一个账号同时只能有一个在线
	s.TryKickOff(svc.ClientId, svc.UserId)

	//发送离线消息
	s.CheckOfflineMsg(svc.UserId)

	gateIp, err := util.LocalIP()
	if nil != err {
		log.Error(err)
		//TODO 关闭链接？还是重试？
		s.DelService(svc)
		return
	}

	onlineReq := &meta.DataOnlineRequest{}
	onlineReq.ClientId = svc.ClientId
	onlineReq.UserId = svc.UserId
	onlineReq.GateIp = gateIp
	onlineReq.GatePort = config.RpcServicePort
	onlineReq.Platform = svc.Platform

	//
	_, err = dataCli.Online(onlineReq)
	if nil != err {
		log.Error(err)
		s.DelService(svc)
		return
	}
}

func (s *Server) CheckOfflineMsg(userId string) {
	resp, err := dataCli.GetOfflineMsgs(userId)
	if nil != err {
		log.Error(err)
		return
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
		log.Debugln("delete service ", svc)
		delete(s.Services, svc.UserId)
	}
}

/*
因为需要辨别每次的链接是否是同一个手机，所以需要手机根据硬件生成一个唯一标识来当作clientId,
clientId格式：OS系统手机硬件唯一对应标识．例如：ios123abb
*/
func (s *Server) ParseClientId(clientId string) (string, error) {
	if strings.HasPrefix(clientId, "ios") {
		return "ios", nil
	}

	return "android", nil
	//	log.Errorln("invalid clientId:", clientId)
	//	return "", errors.New("clientId invalid")
}

//TODO
func (s *Server) Auth(userId, token string) error {
	return nil
}

func (s *Server) CronEvery() {
	for {
		time.Sleep(time.Second * time.Duration(s.Keepalive))
		log.Infoln("it seems", len(s.Services), "clients coneccting,begin to scaning alive")
		for _, user := range s.Services {
			log.Debugf("check clientId:%s,userId:%s is alive?", user.ClientId, user.UserId)
			if !user.IsAlive(s.Keepalive) {
				log.Debugf("clientId:%s,userId:%s is not alive", user.ClientId, user.UserId)
				s.DelService(user)

				if nil != user.Conn {
					err := user.Conn.Close()
					if nil != err {
						log.Errorf("close conn for:%+v,clientId:%s,userId:%s,err:%s", *user, user.ClientId, user.UserId, err.Error())
					} else {
						log.Infof("close conn for:%+v,clientId:%s,userId:%s,err:%s", *user, user.ClientId, user.UserId, err.Error())
					}
				}

				//TODO 通知session
			}
		}
	}
}
