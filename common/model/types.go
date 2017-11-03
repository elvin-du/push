package model

type OfflineMsg struct {
	AppID    string
	ClientID string
	PacketID uint16
	Kind     int
	Content  string
	Extra    string
	CreateAt uint64
}

type App struct {
	ID                     string `json:"id"`
	Secret                 string `json:"secret"`
	AuthType               uint16 `json:"auth_type"`
	Name                   string `json:"name"`
	Description            string `json:"description"`
	Status                 byte   `json:"status"`
	CreatedAt              uint64 `json:"created_at"`
	UpdatedAt              uint64 `json:"updated_at"`
	BundleID               string `json:"bundle_id"`
	Cert                   string `json:"cert"`
	CertPassword           string `json:"cert_password"`
	CertProduction         string `json:"cert_production"`
	CertPasswordProduction string `json:"cert_password_production"`
}
