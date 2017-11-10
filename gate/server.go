package main

import (
	"encoding/json"
	"fmt"
	"gokit/log"
	"io"
	"net"
	"os"
	"os/signal"
	"push/common/model"
	gateMsg "push/gate/message"
	"push/gate/mqtt"
	"push/gate/service/config"
	"push/gate/service/session"
	"push/pb"
	"syscall"
	"time"

	"github.com/surgemq/message"
	"google.golang.org/grpc"
)

var (
	defaultServer = &Server{
		UserManager: defaultUserManager,
		Keepalive:   60 * 5, //五分钟检查一次客户端连接情况
	}
)

type Server struct {
	*UserManager
	Keepalive   int64 //单位：秒
	rcpListener net.Listener
	tcpListener net.Listener
}

/*
开始监听RPC端口
*/
func (s *Server) StartRPCServer() {
	var err error = nil
	s.rcpListener, err = net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", config.RPC_SERVICE_PORT))
	if err != nil {
		panic(err)
	}
	defer func() {
		s.rcpListener.Close()
		log.Infof("rpc connection closed")
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		stop := <-ch
		log.Errorf("receive signal '%v'", stop)
		s.Close(stop)
		gateMsg.DefaultMessageManager.Sync()
		os.Exit(1)
	}()

	log.Infof("starting rpc server on %d", config.RPC_SERVICE_PORT)

	srv := grpc.NewServer()

	pb.RegisterGateServer(srv, &Gate{})
	srv.Serve(s.rcpListener)
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

	err := s.authConnection(ses)
	if nil != err {
		log.Errorln(err)
		return
	}

	err = s.Online(ses)
	if nil != err {
		log.Errorln(err)
		return
	}

	//发送离线消息
	s.CheckOfflineMsgs(ses)
}

func (s *Server) authConnection(ses *mqtt.Session) (err error) {
	connAckMsg := message.NewConnackMessage()

	defer func() {
		//回应connect消息
		if nil != err {
			connAckMsg.SetReturnCode(message.ErrNotAuthorized)
			err2 := ses.SendMsg(connAckMsg)
			if nil != err {
				log.Error(err2)
			}

			//直接关闭连接
			ses.Close(err)
		}
	}()

	var connMsg *message.ConnectMessage
	//获取客户端链接信息
	connMsg, err = ses.GetConnectMessage()
	if nil != err {
		log.Error(err)
		return err
	}

	ses.ClientID = string(connMsg.ClientId())
	err = s.ValidateClientID(ses.ClientID)
	if nil != err {
		log.Error(err)
		return err
	}
	ses.AppID = string(connMsg.Username())
	log.Debugf("come to connect,app_id: %s,clientid:%s", ses.AppID, ses.ClientID)

	//合法性检验
	err = s.Auth(ses.AppID, string(connMsg.Password()))
	if nil != err {
		log.Error(err)
		return err
	}

	//连接成功
	connAckMsg.SetReturnCode(message.ConnectionAccepted)
	err = ses.SendMsg(connAckMsg)
	if nil != err {
		log.Error(err)
		return err
	}

	ses.SetTouchTime(time.Now().Unix())
	ses.OnClose(s.OnSessionClose)
	ses.OnReadPacket(OnRead)
	ses.OnSendPacket(OnSend)
	//启动两个goroutine进行读写
	ses.Start()

	return nil
}

func (s *Server) Online(ses *mqtt.Session) error {
	ses2 := &session.Session{
		AppID:          ses.AppID,
		ClientID:       ses.ClientID,
		GateServerIP:   config.SERVER_IP,
		GateServerPort: config.RPC_SERVICE_PORT,
	}

	err := session.Update(ses2)
	if nil != err {
		log.Error(err)
		return err
	}

	s.PutUser(NewUser(ses))
	log.Infof("app_id:%s,client_id:%s online", ses.AppID, ses.ClientID)

	return nil
}

func (s *Server) CheckOfflineMsgs(ses *mqtt.Session) {
	msgs, err := model.OfflineMsgModel().List(ses.AppID, ses.ClientID)
	if nil != err {
		log.Errorln(err)
		return
	}

	log.Debugf("found %d offline msg for app_id:%s,client_id:%s", len(msgs), ses.AppID, ses.ClientID)
	for _, v := range msgs {
		msg := Message{}
		msg.Content = v.Content
		msg.ID = v.ID
		bin, err := json.Marshal(msg)
		if nil != err {
			log.Errorln(err)
			return
		}

		go ses.Push(bin)
	}
}

func (s *Server) PutUser(u *User) {
	s.UserManager.Put(u)
}

func (s *Server) RemoveUser(appID, clientID string) {
	s.UserManager.Remove(appID, clientID)
}

//TODO 要求长度，大小写等
func (s *Server) ValidateClientID(clientId string) error {
	return nil
}

//TODO
func (s *Server) Auth(appID, appSecret string) error {
	return nil
}

func (s *Server) OnSessionClose(ses *mqtt.Session, err error) {
	if io.EOF == err {
		log.Infof("app_id:%s,client_id:%s session close,err:%s", ses.AppID, ses.ClientID, err.Error())
	} else {
		log.Errorf("app_id:%s,client_id:%s session close,err:%s", ses.AppID, ses.ClientID, err.Error())
	}
	log.Infof("remove user(app_id:%s,client_id:%s) session", ses.AppID, ses.ClientID)
	s.Remove(ses.AppID, ses.ClientID)
	gateMsg.DefaultMessageManager.SyncByAccount(ses.Key())
}

func (s *Server) Close(reason interface{}) {
	s.tcpListener.Close()
	s.rcpListener.Close()
}
