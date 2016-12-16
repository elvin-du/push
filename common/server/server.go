package server

import (
	"log"

	"push/common/etcd"
	"push/common/util"
	"time"
)

type Server struct {
	Client            *etcd.ETCDClient
	APPName           string
	ServiceName       string
	ServiceVersion    string
	IP                string
	Port              string
	Meta              map[string]string
	HeartbeatInterval time.Duration
}

func NewServer(app, srvName, srvVer, port string, meta map[string]string, cli *etcd.ETCDClient, heartbeat time.Duration) *Server {
	ip, err := util.LocalIP()
	if nil != err {
		log.Fatal(err)
		//TODO
	}

	return &Server{
		Client:            cli,
		APPName:           app,
		ServiceName:       srvName,
		ServiceVersion:    srvVer,
		IP:                ip,
		Port:              port,
		Meta:              meta,
		HeartbeatInterval: heartbeat,
	}
}

func (s *Server) Start() error {
	err := s.Register()
	if nil != err {
		log.Print(err)
		return err
	}

	go func() error {
		for {
			time.Sleep(s.HeartbeatInterval)
			err = s.Heartbeat()
			if nil != err {
				log.Print(err)
				return err
			}
		}
	}()

	return nil
}

func (s *Server) Register() error {
	return s.doRegister(nil)
}

func (s *Server) doRegister(meta map[string]string) error {
	err := s.Client.Register(s.APPName, s.ServiceName, s.ServiceVersion, s.IP, s.Port, meta)
	if nil != err {
		log.Println(err)
		return err
	}

	return nil
}

func (s *Server) Heartbeat() error {
	meta := map[string]string{"weight": "5"}
	return s.doRegister(meta)
}

func (s *Server) Get() (ip string, port string, err error) {
	return s.Client.Get(s.APPName, s.ServiceName, s.ServiceVersion)
}
