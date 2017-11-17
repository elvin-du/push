//+build TODO

package main

import (
	"gokit/util"
	"push/common/model"
)

func SaveMsg(msg *Message) error {
	message := &model.Message{}
	message.AppID = msg.AppID
	message.RegID = msg.RegID
	message.Content = msg.Content
	message.CreatedAt = util.Timestamp()
	message.Extras = msg.Extras
	message.ID = msg.ID
	message.TTL = msg.TTL

	return model.MessageModel().Insert(message)
}
