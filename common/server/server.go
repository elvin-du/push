package server

import (
	"hscore/log"
	"net"
	"push/common/client/etcd"
	"push/common/util"
	"time"

	"google.golang.org/grpc"
)

type RPCServer struct {
	APPName           string
	ServiceName       string
	ServiceVersion    string
	IP                string
	Port              string
	Meta              map[string]string
	HeartbeatInterval time.Duration
	Server            *grpc.Server
}

func NewRPCServer(app, serviceName, serviceVer, port string, meta map[string]string, heartbeat time.Duration, server *grpc.Server) *RPCServer {
	ip, err := util.LocalIP()
	if nil != err {
		log.Fatal(err)
		//TODO
	}

	return &RPCServer{
		APPName:           app,
		ServiceName:       serviceName,
		ServiceVersion:    serviceVer,
		IP:                ip,
		Port:              port,
		Meta:              meta,
		HeartbeatInterval: heartbeat,
		Server:            server,
	}
}

func (s *RPCServer) Run() {
	l, err := net.Listen("tcp", s.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Infoln("gate rpc listening at:", s.Port)

	go func() {
		for {
			err = s.Heartbeat()
			if nil != err {
				//注册不成功也不返回，只是记录下来
				log.Errorln(err)
				continue
			}

			time.Sleep(s.HeartbeatInterval)
		}
	}()

	err = s.Server.Serve(l)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
}

func (s *RPCServer) Register() error {
	return s.doRegister(nil)
}

func (s *RPCServer) doRegister(meta map[string]string) error {
	client, err := etcd.GetClient()
	if nil != err {
		log.Errorln(err)
		return err
	}

	err = client.Register(s.APPName, s.ServiceName, s.ServiceVersion, s.IP, s.Port, meta)
	if nil != err {
		log.Errorln(err)
		return err
	}

	return nil
}

/*
每次心跳时都会重新注册一次
*/
func (s *RPCServer) Heartbeat() error {
	//TODO 权重用来做负载均衡
	meta := map[string]string{"weight": "5"}
	return s.doRegister(meta)
}

func (s *RPCServer) Get() (ip string, port string, err error) {
	client, err := etcd.GetClient()
	if nil != err {
		log.Errorln(err)
		return "", "", err
	}

	return client.Get(s.APPName, s.ServiceName, s.ServiceVersion)
}
