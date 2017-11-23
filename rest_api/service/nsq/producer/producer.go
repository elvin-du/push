package producer

import (
	"container/list"
	"gokit/log"

	"github.com/nsqio/go-nsq"
)

var (
	CACHE_SIZE = 32
)

type NsqProducer struct {
	Queue      *list.List
	cacheInCh  chan []byte
	cacheOutCh chan []byte
	producer   *nsq.Producer
	topic      string
}

func NewNsqProducer(addr, topic string) *NsqProducer {
	producer, err := nsq.NewProducer(addr, nsq.NewConfig())
	if nil != err {
		log.Fatal(err) //直接崩溃
	}

	return &NsqProducer{
		Queue:      list.New(),
		cacheInCh:  make(chan []byte, CACHE_SIZE),
		cacheOutCh: make(chan []byte, CACHE_SIZE),
		producer:   producer,
		topic:      topic,
	}
}

func (np *NsqProducer) Start() {
	go np.ReadLoop()
	go np.WriteLoop()
}

func (np *NsqProducer) ReadLoop() {
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

func (np *NsqProducer) WriteLoop() {
	for {
		select {
		case msg := <-np.cacheOutCh:
			err := np.publish(np.topic, msg)
			if nil != err {
				log.Errorf("message:%s into nsq failed,err:%s", string(msg), err.Error())
				continue
			}
			log.Debugf("message:%s into nsq succcess", string(msg))
		}
	}
}

func (np *NsqProducer) publish(topic string, data []byte) error {
	return np.producer.Publish(topic, data)
}

func (np *NsqProducer) Publish(data []byte) error {
	log.Infof("message:%s into out queue", string(data))
	np.cacheInCh <- data
	return nil
}
