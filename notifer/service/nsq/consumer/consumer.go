package consumer

import (
	"gokit/log"
	"push/notifer/service/config"

	"github.com/nsqio/go-nsq"
)

var (
	NSQ_TOPIC_SINGLE      = "topic_push_single"
	NSQ_CHANNEL_SINGLE    = "channel_push_single"
	NSQ_TOPIC_BROADCAST   = "topic_push_broadcast"
	NSQ_CHANNEL_BROADCAST = "channel_push_broadcast"
)

type MsgHandler interface {
	Process(msg interface{}) error
}

type nsqConsumer struct {
	singleConsumer      *nsq.Consumer
	broadcastConsumer   *nsq.Consumer
	singleMsgHandler    MsgHandler
	broadcastMsgHandler MsgHandler
	stopCh              chan byte
}

func NewNsqConsumer(singleMsgHandler MsgHandler, broadcastMsgHandler MsgHandler) *nsqConsumer {
	return &nsqConsumer{
		singleMsgHandler:    singleMsgHandler,
		broadcastMsgHandler: broadcastMsgHandler,
	}
}

func (n *nsqConsumer) Run() {
	var err error = nil
	cfg := nsq.NewConfig()

	n.singleConsumer, err = nsq.NewConsumer(NSQ_TOPIC_SINGLE, NSQ_CHANNEL_SINGLE, cfg)
	if nil != err {
		log.Fatalln(err)
	}
	n.singleConsumer.AddHandler(nsq.HandlerFunc(func(msg *nsq.Message) error {
		return n.singleMsgHandler.Process(msg)
	}))

	err = n.singleConsumer.ConnectToNSQLookupds(config.NSQ_LOOKUPD_ADDRS)
	if nil != err {
		log.Fatalln(err)
	}

	//广播
	n.broadcastConsumer, err = nsq.NewConsumer(NSQ_TOPIC_BROADCAST, NSQ_CHANNEL_BROADCAST, cfg)
	if nil != err {
		log.Fatalln(err)
	}
	n.broadcastConsumer.AddHandler(nsq.HandlerFunc(func(msg *nsq.Message) error {
		return n.broadcastMsgHandler.Process(msg)
	}))

	n.broadcastConsumer.ConnectToNSQLookupds(config.NSQ_LOOKUPD_ADDRS)
	if nil != err {
		log.Fatalln(err)
	}

	<-n.stopCh
}

func (n *nsqConsumer) Stop() {
	n.stopCh <- 1
}
