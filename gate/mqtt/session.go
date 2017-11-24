package mqtt

import (
	"errors"
	"gokit/log"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/surgemq/message"
	"gopkg.in/mgo.v2/bson"
)

var (
	DEFAULT_OUT_CH_SIZE   = 1024
	DEFAULT_WRITE_TIMEOUT = time.Duration(10)  //sec
	KEEPALIVE             = time.Duration(120) //sec
)

// session close flag
const (
	sessionFlagClosed = iota
	sessionFlagOpen
)

type Session struct {
	ID           string
	Conn         net.Conn
	writeTimeout time.Duration
	touchTime    int64
	outCh        chan []byte

	//close
	//	closeChan   chan byte // close write loop
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
		ID:           bson.NewObjectId().Hex(),
		Conn:         conn,
		writeTimeout: DEFAULT_WRITE_TIMEOUT,
		touchTime:    time.Now().Unix(),
		outCh:        make(chan []byte, DEFAULT_OUT_CH_SIZE),
		//		closeChan:    make(chan byte),
		closeWait: new(sync.WaitGroup),
		closeFlag: sessionFlagClosed,
	}

	ses.Conn.SetWriteDeadline(time.Now().Add(ses.writeTimeout * time.Second))

	return ses
}

func (s *Session) Start() {
	if atomic.CompareAndSwapInt32(&s.closeFlag, sessionFlagClosed, sessionFlagOpen) {
		s.closeWait.Add(1)
		go s.ReadLoop()

		s.closeWait.Add(1)
		go s.WriteLoop()

		go s.CronEvery()
	} else {
		log.Errorln("session status unnormal")
	}
}

func (s *Session) Close(reason error) {
	if atomic.CompareAndSwapInt32(&s.closeFlag, sessionFlagOpen, sessionFlagClosed) {
		if io.EOF == reason {
			log.Infoln(reason)
		} else {
			log.Errorln(reason)
		}

		s.closeReason = reason
		//		close(s.closeChan)
		s.Conn.Close()
		s.closeWait.Wait()
		s.closeCallback(s, reason)
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
			err := errors.New("session_id:" + s.ID + " no beatheart,so close socket")
			s.Close(err)
			return
		}

		if s.IsClosed() {
			return
		}
	}
}
