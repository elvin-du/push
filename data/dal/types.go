package dal

type ClientInfo struct {
	Id           string `json:"id"`
	GateServerIp string `json:"gate_server_ip"`
	UserId       string `json:"user_id"`
	Platform     string `json:"platform"`
	Status       byte   `json:"status"`
	CreatedAt    uint64 `json:"created_at"`
	UpdatedAt    uint64 `json:"updated_at"`
}

type OfflineMsg struct {
	Id        string `json:"id"`
	ClientId  string `json:"client_id"`
	UserId    string `json:"user_id"`
	Kind      int    `json:"kind"`
	Content   string `json:"content"`
	Extra     string `json:"extra"`
	CreatedAt uint64 `json:"created_at"`
}
