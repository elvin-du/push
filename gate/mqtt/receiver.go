package mqtt

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"

	"github.com/surgemq/message"
)

var (
	E_WRITE_ERROR = errors.New("write error")
)

func (s *Service) WriteLoop() error {
	for {
		select {
		case out := <-s.outCh:
			n, err := s.Conn.Write(out)
			if nil != err {
				log.Println(err)
				return err
			}
			log.Println("wrote number:", n)
		}
	}
	return E_WRITE_ERROR
}

// Read a raw message from conn
func (s *Service) readRaw() ([]byte, error) {
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
		log.Println(err)
		return nil, err
	}

	msg := message.NewConnectMessage()
	_, err = msg.Decode(buf)
	if nil != err {
		log.Println(err)
		return nil, err
	}

	return msg, err
}
