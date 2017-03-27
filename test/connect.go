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
	conn, err = net.Dial("tcp", ":60001")
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

	err = connMsg.SetClientId([]byte("ios123456789123456789123456789123456789123456789"))
	if nil != err {
		log.Println(err)
		return err
	}

	connMsg.SetCleanSession(false)
	err = Send(connMsg)
	if nil != err {
		log.Println(err)
		return err
	}

	return nil
}
