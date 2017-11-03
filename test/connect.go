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
	connMsg.SetUsername([]byte("87c154323ef0d204"))                 //api_id
	connMsg.SetPassword([]byte("ba8ed065e670d0118261579fd3c1fd52")) //api_secret
	err = Send(connMsg)
	if nil != err {
		log.Println(err)
		return err
	}

	return nil
}
