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

	err = connMsg.SetClientId([]byte("IOS123"))
	if nil != err {
		log.Println(err)
		return err
	}

	connMsg.SetCleanSession(false)
	connMsg.SetUsername([]byte("appid123"))
	connMsg.SetPassword([]byte("app_secret123"))
	err = Send(connMsg)
	if nil != err {
		log.Println(err)
		return err
	}

	return nil
}
