package main

import (
	"gokit/util"
	"push/common/model"
	"time"
)

func SaveMsg(msg *Message) error {
	message := &model.Message{}
	message.AppID = msg.AppID
	message.RegID = msg.RegID
	message.Content = msg.Content
	message.CreatedAt = util.Timestamp()
	message.Extras = msg.Extras
	message.ID = msg.ID
	message.TTL = uint64(time.Now().Add(time.Second * time.Duration(msg.TTL)).Unix())

	return model.MessageModel().Insert(message)
}
