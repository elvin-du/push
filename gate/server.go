package main

import (
	"flag"
	"fmt"
	"net"
	//	"push/common/server"
	//	"push/common/util"
	//	. "push/errors"
	"os"
	"os/signal"
	"push/common/grpclb"
	"push/gate/model"
	"push/gate/mqtt"
	"push/gate/service/config"
	mylog "push/gate/service/log"
	"push/gate/service/session"
	"push/pb"
	"strings"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	log "github.com/sirupsen/logrus"
	"github.com/surgemq/message"
	"google.golang.org/grpc"
)

const (
	RPC_SERVICE_NAME = "gate"
)

var (
	defaultServer = &Server{
		Services:  make(map[string]*mqtt.Service),
		Keepalive: 60 * 5, //五分钟检查一次客户端连接情况
	}
)

type Server struct {
	Services  map[string]*mqtt.Service //key:AppName+clientId
	Keepalive int64                    //单位：秒
}

var (
	reg = flag.String("reg", "http://127.0.0.1:2379", "register etcd address")
)

/*
开始监听RPC端口
*/
func (s *Server) StartRPCServer() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", config.RPC_SERVICE_PORT))
	if err != nil {
		panic(err)
	}

	//TODO
	err = grpclb.Register(RPC_SERVICE_NAME, config.SERVER_IP, config.RPC_SERVICE_PORT, *reg, time.Second*10, 15)
	if err != nil {
		panic(err)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		stop := <-ch
		log.Printf("receive signal '%v'", stop)
		grpclb.UnRegister()
		os.Exit(1)
	}()

	log.Printf("starting hello service at %d", config.RPC_SERVICE_PORT)

	grpc_logrus.ReplaceGrpcLogger(mylog.LogrusEntry)
	srv := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(
				grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor),
			),
			grpc_logrus.UnaryServerInterceptor(mylog.LogrusEntry),
		))

	pb.RegisterGateServer(srv, &Gate{})
	srv.Serve(lis)
}

//开始监听客户端的连接
func (s *Server) StartTcpServer() {
	l, err := net.Listen("tcp", ":"+config.TCP_PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer l.Close()

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

	err = s.Online(svc)
	if nil != err {
		log.Errorln(err)
		return
	}

	//发送离线消息
	s.CheckOfflineMsg(svc.ClientId)
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
	log.Debugln("come to connect,clientid:", string(connMsg.ClientId()))

	svc.ClientId = string(connMsg.ClientId())

	//合法性检验
	err = s.Auth(svc.AppName, string(connMsg.Password()))
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

func (s *Server) Online(svc *mqtt.Service) error {
	var ses session.Session
	ses.ClientID = svc.ClientId
	ses.Platform = svc.Platform
	ses.GateServerIP = config.SERVER_IP
	ses.GateServerPort = config.RPC_SERVICE_PORT

	err := ses.Save()
	if nil != err {
		log.Error(err)
		s.DelService(svc)
		return err
	}

	return nil
}

func (s *Server) CheckOfflineMsg(clientId string) {
	msgs, err := model.OfflineMsgModel().Get(clientId)
	if nil != err {
		log.Errorln(err)
		return
	}

	log.Debugf("found %d offline msg for clientId:%s", len(msgs), clientId)
	svc := s.Services[clientId]
	for _, v2 := range msgs {
		go svc.Push(uint16(v2.PacketID), []byte(v2.Content))
	}
}

func (s *Server) SetService(svc *mqtt.Service) {
	s.Services[svc.AppName+svc.ClientId] = svc
}

func (s *Server) DelService(svc *mqtt.Service) {
	if nil == svc {
		return
	}

	delete(s.Services, svc.AppName+svc.ClientId)

	//	services := make([]*mqtt.Service, 0, 0)

	//	svcs := s.Services[svc.UserId]
	//	for _, service := range svcs {
	//		if service.UserId == svc.UserId &&
	//			service.ClientId == svc.ClientId &&
	//			service.Platform == svc.Platform {
	//			log.Debugln("delete service ", svc)
	//			continue
	//		}
	//		services = append(services, service)
	//	}

	//	s.Services[svc.UserId] = services
}

/*
因为需要辨别每次的链接是否是同一个手机，所以需要手机根据硬件生成一个唯一标识来当作clientId,
clientId格式：OS系统手机硬件唯一对应标识．例如：IOS123abb
*/
func (s *Server) ParseClientId(clientId string) (string, error) {
	if strings.HasPrefix(clientId, "IOS") {
		return "IOS", nil
	}

	return "ANDROID", nil
	//	log.Errorln("invalid clientId:", clientId)
	//	return "", errors.New("clientId invalid")
}

//TODO
func (s *Server) Auth(AppName, appSecret string) error {
	return nil
}

func (s *Server) CronEvery() {
	for {
		time.Sleep(time.Second * time.Duration(s.Keepalive))
		log.Infoln("it seems", len(s.Services), "clients coneccting,begin to scaning alive")

		for _, user := range s.Services {
			log.Debugf("check clientId:%s,AppName:%s is alive?", user.ClientId, user.AppName)
			if !user.IsAlive(s.Keepalive) {
				log.Debugf("clientId:%s,AppName:%s is not alive", user.ClientId, user.AppName)
				s.DelService(user)

				if nil != user.Conn {
					err := user.Conn.Close()
					if nil != err {
						log.Errorf("close conn for:%+v,clientId:%s,AppName:%s,err:%s", *user, user.ClientId, user.AppName, err.Error())
					} else {
						log.Infof("close conn for:%+v,clientId:%s,AppName:%s,err:%s", *user, user.ClientId, user.AppName, err.Error())
					}
				}

				//TODO 通知session
			}
		}
	}
}
