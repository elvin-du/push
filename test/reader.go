package main

import (
	"encoding/binary"
	"log"

	"github.com/surgemq/message"
)

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
