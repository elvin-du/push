package mqtt

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/surgemq/message"
)

var (
	DEFAULT_OUT_CH_SIZE   = 1024
	DEFAULT_WRITE_TIMEOUT = time.Duration(10) //sec
	KEEPALIVE             = time.Duration(60) //sec
)

// session close flag
const (
	sessionFlagClosed = iota
	sessionFlagOpen
)

type Session struct {
	AppID        string
	ClientID     string
	Platform     string
	Conn         net.Conn
	writeTimeout time.Duration
	touchTime    int64
	outCh        chan []byte

	//close
	closeChan   chan byte // close write loop
	closeWait   *sync.WaitGroup
	closeFlag   int32
	closeReason error

	// callback function
	sendPacketCallback func(*Session, []byte) error          // 发送Packet的回调
	readPacketCallback func(*Session, message.Message) error // 接收Packet的回调
	closeCallback      func(*Session, error)
}

func NewSession(conn net.Conn) *Session {
	ses := &Session{
		Conn:         conn,
		writeTimeout: DEFAULT_WRITE_TIMEOUT,
		touchTime:    time.Now().Unix(),
		outCh:        make(chan []byte, DEFAULT_OUT_CH_SIZE),
		closeChan:    make(chan byte),
		closeWait:    new(sync.WaitGroup),
		closeFlag:    sessionFlagClosed,
	}

	ses.Conn.SetWriteDeadline(time.Now().Add(ses.writeTimeout * time.Second))

	return ses
}

func (s *Session) Key() string {
	return fmt.Sprintf("%s+%s", s.AppID, s.ClientID)
}

func (s *Session) Start() {
	s.closeWait.Add(1)
	go s.ReadLoop()

	s.closeWait.Add(1)
	go s.WriteLoop()

	go s.CronEvery()
}

func (s *Session) Close(reason error) {
	if atomic.CompareAndSwapInt32(&s.closeFlag, sessionFlagOpen, sessionFlagClosed) {
		fmt.Println(reason)
		s.closeReason = reason
		close(s.closeChan)
		s.closeCallback(s, reason)
	} else {
		//TODO
	}
}

func (s *Session) SetWriteTimeout(d time.Duration) {
	s.writeTimeout = d
}

func (s *Session) SetTouchTime(t int64) {
	s.touchTime = t
}

// check session is closed or not.
func (s *Session) IsClosed() bool {
	return atomic.LoadInt32(&s.closeFlag) == sessionFlagClosed
}

// 注册packet read callback.
func (s *Session) OnReadPacket(callback func(*Session, message.Message) error) {
	s.readPacketCallback = callback
}

// 注册session send callback.
func (s *Session) OnSendPacket(callback func(*Session, []byte) error) {
	s.sendPacketCallback = callback
}

// 注册session close callback.
func (s *Session) OnClose(callback func(*Session, error)) {
	s.closeCallback = callback
}

func (s *Session) IsAlive() bool {
	return time.Now().Unix()-s.touchTime < int64(KEEPALIVE)
}

func (s *Session) CronEvery() {
	for {
		time.Sleep(KEEPALIVE * time.Second)
		if !s.IsAlive() {
			//TODO maybe should ping first
			err := errors.New("No Beatheart")
			s.Close(err)
		}
	}
}
