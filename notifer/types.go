package main

import (
	"fmt"
	"push/common/types"
)

type Message struct {
	*types.Message
}

func NewMessage() *Message {
	return &Message{
		&types.Message{},
	}
}

func (m *Message) Key() string {
	return fmt.Sprintf("%s:%s", m.AppID, m.RegID)
}

type session struct {
	AppID          string `json:"app_id" redis:"app_id"`
	RegID          string `json:"reg_id" redis:"reg_id"`
	GateServerIP   string `json:"gate_server_ip" redis:"gate_server_ip"`
	GateServerPort string `json:"gate_server_port" redis:"gate_server_port"`
}
