package model

import (
	"fmt"
)

type Message struct {
	ID        string                 `json:"id" bson:"_id"`
	AppID     string                 `json:"app_id" bson:"app_id"`
	RegID     string                 `json:"reg_id" bson:"reg_id"`
	Content   string                 `json:"content"`
	TTL       uint64                 `json:"ttl" bson:"ttl"`
	Extras    map[string]interface{} `json:"extras"`
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

//func (a *App) ToMap() map[string]interface{} {
//	m := make(map[string]interface{})
//	m["id"] = a.ID
//	m["secret"] = a.Secret
//	m["auth_type"] = a.AuthType
//	m["name"] = a.Name
//	m["description"] = a.Description
//	m["status"] = a.Status
//	m["created_at"] = a.CreatedAt
//	m["updated_at"] = a.UpdatedAt
//	m["bundle_id"] = a.BundleID
//	m["cert"] = a.Cert
//	m["cert_password"] = a.CertPassword
//	m["cert_production"] = a.CertProduction
//	m["cert_password_production"] = a.CertPasswordProduction

//	return m
//}

type Registry struct {
	ID        string `json:"id" bson:"_id"`
	AppID     string `json:"app_id" bson:"app_id"`
	DevToken  string `json:"dev_token" bson:"dev_token"`
	Platform  string `json:"platform" bson:"platform"`
	CreatedAt uint64 `json:"created_at" bson:"created_at"`
}
