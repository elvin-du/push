package model

import (
	"gokit/log"
	"gokit/util"
	libdb "push/common/db"

	"gopkg.in/mgo.v2/bson"
)

type offlineMsg struct{}

var _offlineMsg *offlineMsg

func OfflineMsgModel() *offlineMsg {
	return _offlineMsg
}

func (om *offlineMsg) List(appID, clientID string) ([]*OfflineMsg, error) {
	c := libdb.MainMgoDB().C("offline.msgs")
	msgs := []*OfflineMsg{}
	err := c.Find(bson.M{"app_id": appID, "client_id": clientID}).All(&msgs)
	if nil != err {
		log.Errorln(err)
		return nil, err
	}

	return msgs, nil
}

func (om *offlineMsg) Insert(msg *OfflineMsg) error {
	c := libdb.MainMgoDB().C("offline.msgs")
	msg.CreateAt = util.Timestamp()
	err := c.Insert(msg)
	if nil != err {
		log.Errorln(err)
		return err
	}

	return nil
}

func (om *offlineMsg) Delete(msgID string) error {
	c := libdb.MainMgoDB().C("offline.msgs")
	err := c.RemoveId(msgID)
	if nil != err {
		log.Errorln(err)
		return err
	}

	return nil
}
