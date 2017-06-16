package main

type Message struct {
	AppName  string `json:"app_name"`
	ClientId string `json:"client_id"`
	Content  string `json:"content"`
	Kind     int32  `json:"kind"`
	Extra    string `json:"extra"`
}
