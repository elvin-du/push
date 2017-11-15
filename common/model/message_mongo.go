package model

import (
	"gokit/log"
	"gokit/util"
	libdb "push/common/db"

	"gopkg.in/mgo.v2/bson"
)

type message struct{}

var _message *message

func MessageModel() *message {
	return _message
}

func (m *message) List(appID, regID string) ([]*Message, error) {
	c := libdb.MainMgoDB().C("messages")
	msgs := []*Message{}
	err := c.Find(bson.M{"app_id": appID, "reg_id": regID, "status": 1}).All(&msgs)
	if nil != err {
		log.Errorln(err)
		return nil, err
	}

	return msgs, nil
}

func (m *message) Insert(msg *Message) error {
	c := libdb.MainMgoDB().C("messages")
	msg.CreatedAt = msg.CreatedAt
	msg.UpdatedAt = msg.UpdatedAt
	err := c.Insert(msg)
	if nil != err {
		log.Errorln(err)
		return err
	}

	return nil
}

//设置已读
func (m *message) Delete(msgID string) error {
	c := libdb.MainMgoDB().C("messages")
	err := c.Update(bson.M{"_id": msgID, "status": 1}, bson.M{"$set": bson.M{"status": 2, "updated_at": util.Timestamp()}})
	if nil != err {
		log.Errorln(err)
		return err
	}

	return nil
}
