package mqtt

import (
	"encoding/binary"
	"errors"
	"fmt"
	"hscore/log"
	"time"

	"github.com/surgemq/message"
)

var (
	E_WRITE_ERROR = errors.New("write error")
)

func (s *Service) ReadLoop() error {
	for {
		msg, _, _, err := s.ReadMessage()
		if nil != err {
			log.Error(err)
			continue
		}

		err = s.Process(msg)
		if nil != err {
			log.Errorln(err)
			continue
		}
	}

	return E_READ_ERROR
}

// ReadPacket read one packet from conn
func (s *Service) ReadMessage() (message.Message, []byte, int, error) {
	if s.readTimeout > 0 {
		s.Conn.SetReadDeadline(time.Now().Add(s.readTimeout))
	} else {
		s.Conn.SetReadDeadline(time.Time{})
	}

	var (
		// buf for head
		b = make([]byte, 5)

		// total bytes read
		n = 0
	)

	for {
		_, err := s.Conn.Read(b[n : n+1])
		if err != nil {
			return nil, b, 0, err
		}

		// 第一个字节是packet标志位，第二个字节开始为packet body的长度编码，采用的是变长编码
		// 在变长编码中，编码的第二个字节开始为0x80时，表示后面还有字节
		if n >= 1 && b[n] < 0x80 {
			break
		}
		n++

	}

	// fmt.Println("[DEBUG] [ReadPacket] Start -", b)

	// 获取剩余长度
	remLen, _ := binary.Uvarint(b[1 : n+1])
	mtype := message.MessageType(b[0] >> 4)

	buf := make([]byte, n+1+int(remLen))
	copy(buf, b[:n+1])

	if remLen == 0 {
		msg, err := mtype.New()
		dn, err := msg.Decode(buf)
		if err != nil {
			return nil, buf, 0, err
		}

		return msg, nil, dn, nil
	}

	_, err := s.Conn.Read(buf[n+1:]) //[len(b)+1:]
	if err != nil {
		return nil, buf, 0, err
	}

	msg, err := mtype.New()
	dn, err := msg.Decode(buf)
	if err != nil {
		log.Error(err)
		return nil, buf, 0, err
	}

	return msg, nil, dn, nil
}

// Read a raw message from conn
func (s *Service) readRaw() ([]byte, error) {
	if s.readTimeout > 0 {
		s.Conn.SetReadDeadline(time.Now().Add(s.readTimeout))
	} else {
		s.Conn.SetReadDeadline(time.Time{})
	}

	var (
		// the message buffer
		buf []byte

		// tmp buffer to read a single byte
		b = make([]byte, 1)

		// total bytes read
		l = 0
	)

	// Let's read enough bytes to get the message header (msg type, remaining length)
	for {
		// If we have read 5 bytes and still not done, then there's a problem.
		if l > 5 {
			return nil, fmt.Errorf("connect/getMessage: 4th byte of remaining length has continuation bit set")
		}

		n, err := s.Conn.Read(b[0:])
		if err != nil {
			//glog.Debugf("Read error: %v", err)
			return nil, err
		}

		// Technically i don't think we will ever get here
		if n == 0 {
			continue
		}

		buf = append(buf, b...)
		l += n

		// Check the remlen byte (1+) to see if the continuation bit is set. If so,
		// increment cnt and continue reading. Otherwise break.
		if l > 1 && b[0] < 0x80 {
			break
		}
	}

	// Get the remaining length of the message
	remlen, _ := binary.Uvarint(buf[1:])
	buf = append(buf, make([]byte, remlen)...)

	for l < len(buf) {
		n, err := s.Conn.Read(buf[l:])
		if err != nil {
			return nil, err
		}
		l += n
	}

	return buf, nil
}

func (s *Service) GetConnectMessage() (*message.ConnectMessage, error) {
	buf, err := s.readRaw()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	msg := message.NewConnectMessage()
	_, err = msg.Decode(buf)
	if nil != err {
		log.Error(err)
		return nil, err
	}

	return msg, err
}
