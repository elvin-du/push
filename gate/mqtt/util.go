package mqtt

import (
	//	"encoding/binary"
	"gokit/log"

	"github.com/surgemq/message"
)

func NewMessage(buf []byte) (message.Message, error) {
	mtype := message.MessageType(buf[0] >> 4)
	msg, err := mtype.New()
	if nil != err {
		log.Error(err)
		return nil, err
	}

	_, err = msg.Decode(buf)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return msg, nil
}
