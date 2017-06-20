package session

import (
	"push/gate/service/db"
)

func Get(clientId string) (*Session, error) {
	var ses Session
	err := db.Redis().HMGET(clientId, []interface{}{"client_id", "platform", "gate_server_ip", "gate_server_port"}, &ses)
	if nil != err {
		return nil, err
	}

	return &ses, nil
}
