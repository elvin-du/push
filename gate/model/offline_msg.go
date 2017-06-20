package model

import (
	"hscore/log"
	"push/gate/service/db"
)

type offlineMsg struct{}

var _offlineMsg *offlineMsg

func OfflineMsgModel() *offlineMsg {
	return _offlineMsg
}

func (om *offlineMsg) Get(clientId string) ([]*OfflineMsg, error) {
	var omsgs = []*OfflineMsg{}
	if err := db.Mysql().Where("client_id=?", clientId).Find(&omsgs).Error; nil != err {
		log.Errorln(err)
		return nil, err
	}

	return omsgs, nil
}
