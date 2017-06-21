package main

import (
	"encoding/json"
	"errors"
)

var (
	APP_NAME_INVALID  = errors.New("app_name invalid")
	CLIENT_ID_INVALID = errors.New("client_id invalid")
	COTENT_INVALID    = errors.New("content invalid")
	KIND_INVALID      = errors.New("kind invalid")
)

type Message struct {
//	AppName  string `json:"app_name"`
	ClientId string `json:"client_id"`
	Content  string `json:"content"`
	Kind     int    `json:"kind"`
	Extra    string `json:"extra"`
}

func ValidMessage(bin []byte) error {
	var msg Message
	err := json.Unmarshal(bin, &msg)
	if nil != err {
		return err
	}

//	if "" == msg.AppName {
//		return APP_NAME_INVALID
//	}

	if "" == msg.ClientId {
		return CLIENT_ID_INVALID
	}

	if "" == msg.Content {
		return COTENT_INVALID
	}

	if 0 == msg.Kind {
		return KIND_INVALID
	}

	return nil
}
