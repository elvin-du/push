/*
处理来自客户端的信息
*/

package main

import (
	"log"
)

type Handler struct {
	session *Session
}

func NewHandler(s *Session) *Handler {
	return &Handler{
		session: s,
	}
}

func (h *Handler) Process(packet []byte) error {
	log.Println("Process:", string(packet))
	h.session.sendMessageChan <- packet
	log.Println("process ", string(packet), "done")
	return nil
}
