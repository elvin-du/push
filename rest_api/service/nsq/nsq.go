package nsq

import (
	"push/rest_api/service/config"
	"push/rest_api/service/nsq/producer"
)

var (
	NSQ_TOPIC_SINGLE    = "topic_push_single"
	NSQ_TOPIC_BROADCAST = "topic_push_broadcast"
)

var (
	SingleProducer    = producer.NewNsqProducer(config.NSQ_ADDR, NSQ_TOPIC_SINGLE)
	BroadcastProducer = producer.NewNsqProducer(config.NSQ_ADDR, NSQ_TOPIC_BROADCAST)
)

func Init() {
	SingleProducer.Init()
	BroadcastProducer.Init()
}
