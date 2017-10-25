package session

import (
	"fmt"
	"push/gate/service/db"
)

var (
	TTL int = 300 //默认５分钟
)

func Start() {

	//TODO load from config file
}

type Session struct {
	AppID          string `json:"app_id"`
	ClientID       string `json:"client_id"`
	Platform       string `json:"platform"`
	GateServerIP   string `json:"gate_server_ip"`
	GateServerPort int    `json:"gate_server_port"`
}

func (s *Session) ToMap() map[string]interface{} {
	m := make(map[string]interface{}, 4)
	m["app_id"] = s.AppID
	m["client_id"] = s.ClientID
	m["platform"] = s.Platform
	m["gate_server_ip"] = s.GateServerIP
	m["gate_server_port"] = s.GateServerPort

	return m
}

//每次保存一次，会自动更新过期时间
func (s *Session) Save() error {
	return db.Redis().HMSETAndEXPIRE(s.RedisKey(), s.ToMap(), TTL)
}

func (s *Session) RedisKey() string {
	return fmt.Sprintf("%s+%s", s.AppID, s.ClientID)
}

func (s *Session) Touch() error {
	return db.Redis().EXPIRE(s.ClientID, TTL)
}
