package main

type Message struct {
	AppID  string `json:"app_id"`
	ClientId string `json:"client_id"`
	Content  string `json:"content"`
	Kind     int32  `json:"kind"`
	Extra    string `json:"extra"`
}

type session struct {
	ClientID       string `json:"client_id" redis:"client_id"`
	Platform       string `json:"platform" redis:"platform"`
	GateServerIP   string `json:"gate_server_ip" redis:"gate_server_ip"`
	GateServerPort string `json:"gate_server_port" redis:"gate_server_port"`
}
