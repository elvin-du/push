package main

import (
	"log"
	"net"
	"time"
)

var (
	SEND_CHAN_SIZE int = 500
)

type Session struct {
	netConn         net.Conn
	sendMessageChan chan []byte
}

func NewSession(conn net.Conn) *Session {
	return &Session{
		netConn:         conn,
		sendMessageChan: make(chan []byte, SEND_CHAN_SIZE),
	}
}

func (s *Session) Start() {
	//没有timeout
	err := s.netConn.SetDeadline(time.Time{})
	if nil != err {
		log.Fatal(err)
	}

	go s.writeLoop()
	go s.readLoop()
}

func (s *Session) writeLoop() {
	for {
		select {
		case msg := <-s.sendMessageChan:
			log.Println("send ", string(msg))
			n, err := s.netConn.Write(msg)
			if nil != err {
				log.Println(err)
				continue
			}
			log.Println("send ", string(msg), "done", n)
		}
	}
}

func (s *Session) readLoop() {
	for {
		bin := make([]byte, 1024)
		n, err := s.netConn.Read(bin)
		if nil != err {
			log.Println(err)
			//TODO 处理断线
			return
		}
		log.Println("received:", string(bin[:n]))
		err = NewHandler(s).Process(bin[:n])
		if nil != err {
			log.Println(err)
			continue
		}
		//TODO handle msg
	}
}
