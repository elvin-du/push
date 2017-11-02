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
	connMsg.SetUsername([]byte("8c2e1fb321f36094"))                 //api_id
	connMsg.SetPassword([]byte("4591a021d339f04dfeed738451142006")) //api_secret
	err = Send(connMsg)
	if nil != err {
		log.Println(err)
		return err
	}

	return nil
}
