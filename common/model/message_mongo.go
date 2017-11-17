package model

import (
	"gokit/log"
	"gokit/util"
	libdb "push/common/db"
	"time"

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
	err := c.Find(bson.M{"app_id": appID, "reg_id": regID, "status": 0, "ttl": bson.M{"$gt": time.Now().Unix()}}).All(&msgs)
	if nil != err {
		log.Errorln(err)
		return nil, err
	}

	return msgs, nil
}

func (m *message) Insert(msg *Message) error {
	c := libdb.MainMgoDB().C("messages")
	msg.TTL = uint64(time.Now().Add(time.Second * time.Duration(msg.TTL)).Unix())
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
	err := c.Update(bson.M{"_id": msgID, "status": 0}, bson.M{"$set": bson.M{"status": 1, "updated_at": util.Timestamp()}})
	if nil != err {
		log.Errorln(err)
		return err
	}

	return nil
}
