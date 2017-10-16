package publisher

import (
	"container/list"
	"gokit/log"
	. "gokit/rabbitmq/message"
	"time"

	"github.com/streadway/amqp"
)

type Publisher struct {
	rc             <-chan *Message
	wc             chan<- *Message
	addr, exchange string
	pool           int
	waitTime       time.Duration
}

func NewPublisher(addr, exchange string, pool int, waitTime time.Duration) *Publisher {
	return &Publisher{addr: addr, exchange: exchange, pool: pool, waitTime: waitTime}
}

func (pub *Publisher) Start() {
	log.Debugln("Prepare infinite channel")
	pub.newInfiniteChan()

	log.Debugln("Starting publish worker")
	for i := 0; i < pub.pool; i++ {
		go pub.startWorker(i, nil, 0)
	}

	log.Infoln("Publish service started")
}

func (pub *Publisher) Publish(msg *Message) {
	log.Debugln("Publish message", msg)
	pub.wc <- msg
}

func (pub *Publisher) startWorker(tid int, lastM *Message, lastWait time.Duration) {
	log.Debugln("#", tid, "started")

	var (
		m  *Message
		ok bool
	)

	if lastM != nil {
		m = lastM
		log.Debugln("#", tid, "Got last unsent message", m)
	} else {
		log.Infoln("#", tid, "Waiting for message")
		m, ok = <-pub.rc
		if !ok {
			log.Errorln("#", tid, "Channel closed, quit worker now")
			return
		}

		log.Debugln("#", tid, "Received message", m)
	}

	conn, err := amqp.Dial(pub.addr)
	if err != nil {
		log.Errorln("#", tid, err)
		pub.restartWorker(tid, m, lastWait+pub.waitTime)
		return
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Errorln("#", tid, err)
		conn.Close()
		pub.restartWorker(tid, m, lastWait+pub.waitTime)
		return
	}

	err = ch.ExchangeDeclare(pub.exchange, "topic", true, false, false, false, nil)
	if err != nil {
		log.Errorln("#", tid, err)
		conn.Close()
		pub.restartWorker(tid, m, lastWait+pub.waitTime)
		return
	}

	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         []byte(m.Content),
	}

	if m.Expiration != "" {
		msg.Expiration = m.Expiration
	}

	err = ch.Publish(pub.exchange, m.RouteKey, false, false, msg)
	if err != nil {
		log.Errorln("#", tid, "Publish message failed, error:", err)
		conn.Close()
		pub.restartWorker(tid, m, 0) //reset to 0
		return
	}

	for {
		log.Infoln("#", tid, "Waiting for message")
		m, ok = <-pub.rc
		if !ok {
			log.Errorln("#", tid, "Channel closed, Quit now")
			return
		}

		log.Debugln("#", tid, "Received message", m)

		msg = amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(m.Content),
		}

		if m.Expiration != "" {
			msg.Expiration = m.Expiration
		}

		err = ch.Publish(pub.exchange, m.RouteKey, false, false, msg)
		if err != nil {
			log.Errorln("#", tid, "Publish message failed, error:", err)
			conn.Close()
			break
		}
	}

	pub.restartWorker(tid, m, 0) //reset to 0
}

func (pub *Publisher) restartWorker(tid int, lastM *Message, wait time.Duration) {
	if wait > 0 {
		log.Infoln("#", tid, "Wait for", wait, "to restart worker")
		time.Sleep(wait)
	}

	log.Infoln("Restarting worker now")
	pub.startWorker(tid, lastM, wait)
}

func (pub *Publisher) newInfiniteChan() {
	rc := make(chan *Message)
	wc := make(chan *Message, 100)

	go pub.chanLoop(wc, rc)

	pub.rc = rc
	pub.wc = wc
}

func (pub *Publisher) chanLoop(r <-chan *Message, w chan<- *Message) {
	cache := list.New()

	for {
		if cache.Len() > 0 {
			select {
			case d := <-r:
				cache.PushBack(d)
			case w <- cache.Front().Value.(*Message):
				cache.Remove(cache.Front())
			}
		} else {
			select {
			case d := <-r:
				cache.PushBack(d)
			}
		}
	}
}
