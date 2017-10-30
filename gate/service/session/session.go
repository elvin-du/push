package session

import (
	"fmt"
	"push/gate/service/db"
)

var (
	TTL int = 300 //默认５分钟
)

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

func (s *Session) RedisKey() string {
	return RedisKey(s.AppID, s.ClientID)
}

func RedisKey(appID, clientID string) string {
	return fmt.Sprintf("%s+%s", appID, clientID)
}

//每次保存一次，会自动更新过期时间
func Update(s *Session) error {
	r := db.MainRedis()
	defer r.Close()

	return r.HMSETAndEXPIRE(s.RedisKey(), s.ToMap(), TTL)
}

func Touch(appID, clientID string) error {
	r := db.MainRedis()
	defer r.Close()

	return r.EXPIRE(RedisKey(appID, clientID), TTL)
}
