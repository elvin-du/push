package main

import (
	"encoding/json"
	"gokit/log"
	"push/common/model"
)

func ValidateMessage(bin []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(bin, &msg)
	if nil != err {
		log.Errorln(err, string(bin))
		return nil, err
	}

	if "" == msg.RegID {
		return nil, REG_ID_INVALID
	}

	if "" == msg.Content {
		return nil, COTENT_INVALID
	}

	if 0 == msg.Kind {
		return nil, KIND_INVALID
	}

	return &msg, nil
}

func ValidateRegisterReq(bin []byte) (*model.Registry, error) {
	reg := model.Registry{}
	err := json.Unmarshal(bin, &reg)
	if nil != err {
		log.Errorln(err)
		return nil, err
	}

	if "" == reg.AppID {
		return nil, APP_ID_INVALID
	}

	if "" == reg.Kind {
		return nil, KIND_INVALID
	}

	if "ios" == reg.Kind {
		if "" == reg.DevToken {
			return nil, DEV_TOKEN_INVALID
		}
	}

	return &reg, nil
}

func ValidateAppReq(bin []byte) (*model.App, error) {
	app := model.App{}
	err := json.Unmarshal(bin, &app)
	if nil != err {
		log.Errorln(err)
		return nil, err
	}

	return &app, nil
}
