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
	var err error

	defer func() {
		s.Close(err)
	}()

	for {
		select {
		case msg := <-s.sendMessageChan:
			log.Println("send ", string(msg))
			var n int = 0
			n, err = s.netConn.Write(msg)
			if nil != err {
				log.Println(err)
				return
			}

			log.Println("send ", string(msg), "done", n)
		}
	}
}

// 循环的read并解析网络数据
//
// 注意:
// 1. 如果数据读取或者解包失败, 会直接关闭Session
func (s *Session) readLoop() {
	var err error

	defer func() {
		s.Close(err)
	}()

	for {
		bin := make([]byte, 1024)
		var n int = 0
		n, err = s.netConn.Read(bin)
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

func (s *Session) Close(reason error) {
	//TODO 处理reason

	close(s.sendMessageChan)

	err := s.netConn.Close()
	if nil != err {
		log.Println(err)
	}
}
