package main

type Message struct {
	AppId    string `json:"app_id"`
	ClientId string `json:"client_id"`
	Content  string `json:"content"`
	Kind     int32  `json:"kind"`
	Extra    string `json:"extra"`
}
