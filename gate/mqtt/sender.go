package mqtt

import (
	"errors"
	"fmt"
	"gokit/log"
	"runtime/debug"
	"time"

	"github.com/surgemq/message"
)

var (
	E_WRITE_ERROR = errors.New("write error")
)

func (s *Session) WriteLoop() {
	var err error = nil

	defer func() {
		// handle panic
		if r := recover(); r != nil {
			log.Errorf("session(id:%s) WriteLoop recover: %v, DEBUG.STACK=%v", s.ID, r, string(debug.Stack()))
		}

		s.closeWait.Done()
		s.Close(err)
	}()

	for {
		select {
		case out := <-s.outCh:
			err = s.Send(out)
			if nil != err {
				log.Errorln(err)
				//TODO 直接return?
				return
			}

			if nil != s.sendPacketCallback {
				go s.sendPacketCallback(s, out)
			}
			//		case <-s.closeChan:
			//			s.Conn.Close()

			//			s.closeWait.Wait()

			//			// invoke session close callback
			//			if s.closeCallback != nil {
			//				s.closeCallback(s, s.closeReason)
			//			}

			//			return
		}
	}
}

func (s *Session) Push(content []byte) error {
	pubMsg := message.NewPublishMessage()
	pubMsg.SetQoS(message.QosAtLeastOnce)
	pubMsg.SetPacketId(0)
	pubMsg.SetTopic([]byte("*"))
	pubMsg.SetPayload(content)

	return s.WriteMsg(pubMsg)
}

//只是把信息放到队列，等待发送Goroutine读取并发送
func (s *Session) WriteMsg(msg message.Message) error {
	log.Infof("write message type:%s,session_id:%s", msg.Desc(), s.ID)

	buf := make([]byte, msg.Len())
	n, err := msg.Encode(buf)
	if nil != err {
		log.Error(err)
		return err
	}

	//把数据放到队列
	s.outCh <- buf[:n]

	return nil
}

func (s *Session) SendMsg(msg message.Message) error {
	log.Infof("send message:%s,sesssion_id:%s", msg.Desc(), s.ID)

	buf := make([]byte, msg.Len())
	n, err := msg.Encode(buf)
	if nil != err {
		log.Error(err)
		return err
	}

	return s.Send(buf[:n])
}

//真正的把数据发送到客户端
func (s *Session) Send(data []byte) error {
	s.Conn.SetWriteDeadline(time.Now().Add(s.writeTimeout * time.Second))

	n, err := s.Conn.Write(data)
	if nil != err {
		log.Error(err)
		return err
	}

	if n != len(data) {
		err := fmt.Errorf("data size:%d,but only sent %d success", len(data), n)
		log.Errorln(err)
		return err
	}

	return nil
}
