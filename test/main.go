package main

import (
	"encoding/binary"
	"fmt"
	"hscore/log"
	"net"
	"time"

	"github.com/surgemq/message"
	//	"github.com/surgemq/surgemq/service"
)

var (
	conn net.Conn
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func Connect() error {
	var err error = nil
	conn, err = net.Dial("tcp", ":60001")
	if nil != err {
		log.Println(err)
		return err
	}

	connMsg := message.NewConnectMessage()
	connMsg.SetVersion(4)
	connMsg.SetClientId([]byte("clientid123"))
	connMsg.SetCleanSession(false)
	err = Send(connMsg)
	if nil != err {
		log.Println(err)
		return err
	}

	return nil
}

func Ping() {
	for {
		pingMsg := message.NewPingreqMessage()
		err := Send(pingMsg)
		if nil != err {
			log.Println(err)
			return
		}
		time.Sleep(time.Second * 5)
	}
}

func ReadLoop() {
	for {
		msg, _, _, err := ReadMessage()
		if nil != err {
			log.Println(err)
			return
		}

		Process(msg)
	}
}

func main() {
	err := Connect()
	if nil != err {
		log.Println(err)
		return
	}

	go ReadLoop()

	Ping()
	select {}
}

func Send(msg message.Message) error {
	buf := make([]byte, msg.Len())
	n, err := msg.Encode(buf)
	if nil != err && n == len(buf) {
		log.Println(err)
		return err
	}

	n, err = conn.Write(buf)
	if nil != err {
		log.Println(err)
		return err
	}

	if n != len(buf) {
		err = fmt.Errorf("expected len:%d,got:%d", len(buf), n)
		log.Println(err)
		return err
	}

	return nil
}

// ReadPacket read one packet from conn
func ReadMessage() (message.Message, []byte, int, error) {
	var (
		// buf for head
		b = make([]byte, 5)

		// total bytes read
		n = 0
	)

	for {
		_, err := conn.Read(b[n : n+1])
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

	_, err := conn.Read(buf[n+1:]) //[len(b)+1:]
	if err != nil {
		return nil, buf, 0, err
	}

	msg, err := mtype.New()
	dn, err := msg.Decode(buf)
	if err != nil {
		log.Println(err)
		return nil, buf, 0, err
	}

	return msg, nil, dn, nil
}
