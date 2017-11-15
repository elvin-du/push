package model

import (
	"fmt"
)

type Message struct {
	ID        string `json:"id" bson:"_id"`
	AppID     string `json:"app_id" bson:"app_id"`
	RegID     string `json:"reg_id" bson:"reg_id"`
	Kind      int32
	Content   string
	Extra     string
	Status    byte
	CreatedAt uint64 `json:"created_at" bson:"created_at"`
	UpdatedAt uint64 `json:"updated_at " bson:"updated_at "`
}

func (m *Message) Key() string {
	return fmt.Sprintf("%s:%s", m.AppID, m.RegID)
}

type App struct {
	ID                     string `json:"id" bson:"_id"`
	Secret                 string `json:"secret"`
	AuthType               uint16 `json:"auth_type" bson:"auth_type"`
	Name                   string `json:"name"`
	Description            string `json:"description"`
	Status                 byte   `json:"status"`
	CreatedAt              uint64 `json:"created_at" bson:"created_at"`
	UpdatedAt              uint64 `json:"updated_at" bson:"updated_at"`
	BundleID               string `json:"bundle_id" bson:"bundle_id"`
	Cert                   string `json:"cert"`
	CertPassword           string `json:"cert_password" bson:"cert_password"`
	CertProduction         string `json:"cert_production" bson:"cert_production"`
	CertPasswordProduction string `json:"cert_password_production" bson:"cert_password_production"`
}

type Registry struct {
	ID        string `json:"id" bson:"_id"`
	AppID     string `json:"app_id" bson:"app_id"`
	DevToken  string `json:"dev_token" bson:"dev_token"`
	Kind      string
	CreatedAt uint64 `json:"created_at" bson:"created_at"`
}
