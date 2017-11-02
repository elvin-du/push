package main

import (
	"fmt"
)

type Message struct {
	AppID    string `json:"app_id"`
	Platform string `json:"platform"`
	ClientID string `json:"client_id"`
	Content  string `json:"content"`
	Kind     int32  `json:"kind"`
	Extra    string `json:"extra"`
}

func (m *Message) Key() string {
	return fmt.Sprintf("%s+%s", m.AppID, m.ClientID)
}

type session struct {
	AppID          string `json:"app_id" redis:"app_id"`
	ClientID       string `json:"client_id" redis:"client_id"`
	GateServerIP   string `json:"gate_server_ip" redis:"gate_server_ip"`
	GateServerPort string `json:"gate_server_port" redis:"gate_server_port"`
}
