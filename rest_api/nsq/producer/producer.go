package producer

import (
	"container/list"
	"log"

	"github.com/nsqio/go-nsq"
)

var (
	CACHE_SIZE = 32
)

type nsqProducer struct {
	Queue      *list.List
	cacheInCh  chan []byte
	cacheOutCh chan []byte
	producer   *nsq.Producer
	topic      string
}

func NewNsqProducer(addr, topic string) *nsqProducer {
	producer, err := nsq.NewProducer(addr, nsq.NewConfig())
	if nil != err {
		log.Fatal(err) //直接崩溃
	}

	return &nsqProducer{
		Queue:      list.New(),
		cacheInCh:  make(chan []byte, CACHE_SIZE),
		cacheOutCh: make(chan []byte, CACHE_SIZE),
		producer:   producer,
		topic:      topic,
	}
}

func (np *nsqProducer) Start() {
	go np.ReadLoop()
	go np.WriteLoop()
}

func (np *nsqProducer) ReadLoop() {
	for {
		if np.Queue.Len() > 0 {
			select {
			case m := <-np.cacheInCh:
				np.Queue.PushBack(m)
			case np.cacheOutCh <- np.Queue.Front().Value.([]byte):
				np.Queue.Remove(np.Queue.Front())
			}
		} else {
			select {
			case m := <-np.cacheInCh:
				np.Queue.PushBack(m)
			}
		}
	}
}

func (np *nsqProducer) WriteLoop() {
	for {
		select {
		case msg := <-np.cacheOutCh:
			err := np.publish(np.topic, msg)
			if nil != err {
				log.Println(err)
				continue
			}
		}
	}
}

func (np *nsqProducer) publish(topic string, data []byte) error {
	return np.producer.Publish(topic, data)
}

func (np *nsqProducer) Publish(data []byte) error {
	np.cacheInCh <- data
	return nil
}
