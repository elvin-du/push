package main

import (
	"log"
	//	"time"

	"github.com/surgemq/message"
	"github.com/surgemq/surgemq/service"
)

var (
	mqttAddr = "tcp://127.0.0.1:60001"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	mqttCli := &service.Client{}

	msg := message.NewConnectMessage()
	msg.SetVersion(4)
	msg.SetClientId([]byte("clientid123"))
	err := mqttCli.Connect(mqttAddr, msg)
	if nil != err {
		log.Fatal(err)
	}

	//	for {
	//	time.Sleep(time.Second * 60)
	//		err := mqttCli.Ping(OnComplete)
	err = mqttCli.Ping(OnComplete)
	if nil != err {
		log.Println(err)
		return
	}
	//	}

	//	pubMsg := message.NewPublishMessage()
	//	pubMsg.SetPayload([]byte("1234"))
	//	pubMsg.SetPacketId(5)
	//	pubMsg.SetQoS(message.QosAtLeastOnce)
	//	//	err = pubMsg.SetTopic([]byte("/yd/push/user/b"))
	//	//	if nil != err {
	//	//		log.Fatal(err)
	//	//	}

	//	err = mqttCli.Publish(pubMsg, pubOnComplete)
	//	if nil != err {
	//		log.Fatal(err)
	//	}

	select {}
}
func OnComplete(msg, ack message.Message, err error) error {
	log.Println(ack.Name())
	return nil
}

func subOnPublish(msg *message.PublishMessage) error {

	return nil
}
