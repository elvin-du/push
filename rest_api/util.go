package main

import (
	"encoding/json"
	"gokit/log"
	"push/common/model"
	"push/rest_api/client"
)

func ValidateNotification(bin []byte) (*client.Notification, error) {
	var msg client.Notification
	err := json.Unmarshal(bin, &msg)
	if nil != err {
		log.Errorln(err, string(bin))
		return nil, err
	}

	if "" == msg.Alert {
		return nil, REQUEST_DATA_INVALID
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

	if "" == reg.Platform {
		return nil, PLATFORM_INVALID
	}

	if "ios" == reg.Platform {
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
