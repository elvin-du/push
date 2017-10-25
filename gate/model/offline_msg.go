package model

import (
	"gokit/log"
	"push/gate/service/db"
)

type offlineMsg struct{}

var _offlineMsg *offlineMsg

func OfflineMsgModel() *offlineMsg {
	return _offlineMsg
}

func (om *offlineMsg) Get(appID, clientID string) ([]*OfflineMsg, error) {
	var omsgs = []*OfflineMsg{}
	if err := db.Mysql().Where("app_id=? AND client_id=?", appID, clientID).Find(&omsgs).Error; nil != err {
		log.Errorln(err)
		return nil, err
	}

	return omsgs, nil
}
