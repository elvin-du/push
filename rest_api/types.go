package main

import (
	"encoding/json"
	"errors"
)

var (
	APP_ID_INVALID    = errors.New("app_id invalid")
	CLIENT_ID_INVALID = errors.New("client_id invalid")
	COTENT_INVALID    = errors.New("content invalid")
	KIND_INVALID      = errors.New("kind invalid")
)

type Message struct {
	AppId    string `json:"app_id"`
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

	if "" == msg.AppId {
		return APP_ID_INVALID
	}

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
