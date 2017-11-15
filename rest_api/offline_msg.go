package main

import (
	"gokit/util"
	"push/common/model"
)

func SaveMsg(info *Info) error {
	msg := &model.Message{}
	msg.AppID = info.AppID
	msg.RegID = info.RegID
	msg.Content = info.Content
	msg.CreatedAt = util.Timestamp()
	msg.Extra = info.Extra
	msg.ID = info.ID
	msg.Kind = int32(info.Kind)

	return model.MessageModel().Insert(msg)
}
