package main

import (
	"log"
	"net"

	"github.com/surgemq/message"
)

var (
	conn net.Conn
)

func Connect() error {
	var err error = nil
	conn, err = net.Dial("tcp", ":51001")
	if nil != err {
		log.Println(err)
		return err
	}

	connMsg := message.NewConnectMessage()
	err = connMsg.SetVersion(4)
	if nil != err {
		log.Println(err)
		return err
	}

	err = connMsg.SetClientId([]byte("QQQWWWEEERRR"))
	if nil != err {
		log.Println(err)
		return err
	}

	connMsg.SetCleanSession(false)
	connMsg.SetUsername([]byte("63163c7b40f2abee"))                 //api_id
	connMsg.SetPassword([]byte("283abdfc9123987980d8aabaa7108e6c")) //api_secret
	err = Send(connMsg)
	if nil != err {
		log.Println(err)
		return err
	}

	return nil
}
