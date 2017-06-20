package session

type Session struct {
	ClientID       string `json:"client_id" redis:"client_id"`
	Platform       string `json:"platform" redis:"platform"`
	GateServerIP   string `json:"gate_server_ip" redis:"gate_server_ip"`
	GateServerPort string `json:"gate_server_port" redis:"gate_server_port"`
}
