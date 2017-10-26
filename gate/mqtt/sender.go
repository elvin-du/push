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
			log.Errorf("session WriteLoop recover: %v, DEBUG.STACK=%v", r, string(debug.Stack()))
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
				//TODO 直接断开连接？等待心跳机制检查时删除本服务？
				return
			}

			if nil != s.sendPacketCallback {
				s.sendPacketCallback(s, out)
			}
		case <-s.closeChan:
			s.Conn.Close()

			s.closeWait.Wait()

			// invoke session close callback
			if s.closeCallback != nil {
				s.closeCallback(s, s.closeReason)
			}

			return
		}
	}
}

func (s *Session) Push(packetId uint16, content []byte) error {
	pubMsg := message.NewPublishMessage()
	pubMsg.SetQoS(message.QosAtLeastOnce)
	pubMsg.SetPacketId(packetId)
	pubMsg.SetTopic([]byte("*"))
	pubMsg.SetPayload(content)

	return s.WriteMsg(pubMsg)
}

//只是把信息放到队列，等待发送Goroutine读取并发送
func (s *Session) WriteMsg(msg message.Message) error {
	log.Debugln("write:", msg.Desc())

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
	log.Debugln("SendMsg:", msg.Desc())

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

	log.Debugf("Send begin: data %s", string(data))
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
	log.Debugf("Send success: data %s", string(data))
	return nil
}
