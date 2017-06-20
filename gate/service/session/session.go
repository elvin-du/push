package session

import (
	"push/gate/service/db"
	"time"
)

var (
	TTL time.Duration = time.Second * 300 //默认５分钟
)

func Start() {

	//TODO load from config file
}

type Session struct {
	ClientID       string `json:"client_id"`
	Platform       string `json:"platform"`
	GateServerIP   string `json:"gate_server_ip"`
	GateServerPort string `json:"gate_server_port"`
}

func (s *Session) ToMap() map[string]interface{} {
	m := make(map[string]interface{}, 4)
	m["client_id"] = s.ClientID
	m["platform"] = s.Platform
	m["gate_server_ip"] = s.GateServerIP
	m["gate_server_port"] = s.GateServerPort

	return m
}

//每次保存一次，会自动更新过期时间
func (s *Session) Save() error {
	return db.Redis().HMSETAndEXPIRE(s.ClientID, s.ToMap(), TTL)
}

func (s *Session) Touch() error {
	return db.Redis().EXPIRE(s.ClientID, TTL)
}
