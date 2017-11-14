package main

import (
	"log"
	"net"

	"gokit/util"

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

	return nil
}

func SingIn(appID, appSeceret, regID string) error {
	connMsg := message.NewConnectMessage()
	err := connMsg.SetVersion(4)
	if nil != err {
		log.Println(err)
		return err
	}

	//	err = connMsg.SetClientId([]byte("QQQWWWEEERRR"))
	err = connMsg.SetClientId([]byte(regID))
	if nil != err {
		log.Println(err)
		return err
	}

	appSec, err := util.RC4EncryptToBase64("01e9175ca8805cc2137c44eb86184922", []byte(appID+":"+appSeceret))
	if nil != err {
		log.Println(err)
		return err
	}
	connMsg.SetCleanSession(false)
	connMsg.SetUsername([]byte(appID))  //api_id
	connMsg.SetPassword([]byte(appSec)) //api_secret
	//	connMsg.SetUsername([]byte("63163c7b40f2abee"))                 //api_id
	//	connMsg.SetPassword([]byte("283abdfc9123987980d8aabaa7108e6c")) //api_secret

	err = Send(connMsg)
	if nil != err {
		log.Println(err)
		return err
	}

	return nil
}
