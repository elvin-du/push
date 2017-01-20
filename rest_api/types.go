package main

type Message struct {
	UserId  string `json:"user_id"`
	Content string `json:"content"`
	Kind    int    `json:"kind"`
	Extra   string `json:"extra"`
}
