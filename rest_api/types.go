package main

import (
	"encoding/json"
	"errors"
)

const (
	PUSH_PLATFORM_ANDROID = "android"
	PUSH_PLATFORM_IOS     = "ios"
)

var (
	CLIENT_ID_INVALID = errors.New("client_id invalid")
	COTENT_INVALID    = errors.New("content invalid")
	KIND_INVALID      = errors.New("kind invalid")
)

type Message struct {
	ClientID     string `json:"client_id"`
	Platform     string `json:"platform"` //android,ios
	IsProduction bool   `json:"is_production"`
	Content      string `json:"content"`
	Kind         int    `json:"kind"`
	Extra        string `json:"extra"`
}

type Info struct {
	AppID string `json:"app_id"`
	*Message
}

func ValidMessage(bin []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(bin, &msg)
	if nil != err {
		return nil, err
	}

	if "" == msg.ClientID {
		return nil, CLIENT_ID_INVALID
	}

	if "" == msg.Content {
		return nil, COTENT_INVALID
	}

	if 0 == msg.Kind {
		return nil, KIND_INVALID
	}

	return &msg, nil
}
