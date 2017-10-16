package pub

import (
	"gokit/config"
	"gokit/log"
	//	"gokit/rabbitmq/message"
	"fmt"
	"gokit/rabbitmq/publisher"
	"strings"
	"time"
)

//default to 10sec
var WAIT_STEP = time.Second * 10

var globalPub = make(map[string]*publisher.Publisher)

func CreatePublisher(pubName string) {
	var broker string
	err := config.Get(pubName+":broker", &broker)
	if err != nil {
		log.Fatalln(err)
	}

	log.Debug("Detecting broker type")
	if !strings.HasPrefix(broker, "amqp://") {
		log.Fatal("Unsupported broker type")
	}

	var exchange string
	err = config.Get(pubName+":exchange", &exchange)
	if err != nil {
		log.Fatalln(err)
	}

	var pool int
	err = config.Get(pubName+":pool", &pool)
	if err != nil {
		log.Fatalln(err)
	}

	if pool < 1 {
		pool = 10
	}

	var wait_time int
	err = config.Get(pubName+":reconnect", &wait_time)
	if err != nil {
		log.Fatalln(err)
	}

	if wait_time > 0 {
		WAIT_STEP = time.Duration(wait_time) * time.Second
	}

	pub := publisher.NewPublisher(broker, exchange, pool, WAIT_STEP)
	pub.Start()
	globalPub[pubName] = pub
}

func GetPublisher(pubName string) (*publisher.Publisher, error) {
	pub := globalPub[pubName]
	if nil == pub {
		return nil, fmt.Errorf("Not found %s publisher", pubName)
	}

	return pub, nil
}

//func Publish(message *message.Message) {
//	pub.Publish(message)
//}
