package message

import (
	"gokit/log"
	"push/common/model"
	"sync"
	"time"
)

const (
	ASYNC_INTERNAL = time.Second * 60
)

var DefaultMessageManager = NewMessageManager()

type MessageManager struct {
	lock *sync.RWMutex
	Msgs map[string]*model.Message //key = msg_id
}

func NewMessageManager() *MessageManager {
	manager := &MessageManager{
		lock: new(sync.RWMutex),
		Msgs: make(map[string]*model.Message),
	}
	go manager.AsyncLoop()
	return manager
}

func (mm *MessageManager) IsExist(msgID string) bool {
	mm.lock.RLock()
	defer mm.lock.RUnlock()
	_, ok := mm.Msgs[msgID]
	return ok
}

func (mm *MessageManager) Put(msg *model.Message) {
	mm.lock.Lock()
	defer mm.lock.Unlock()

	mm.Msgs[msg.ID] = msg
}

func (mm *MessageManager) Delete(msgID string) {
	mm.lock.Lock()
	defer mm.lock.Unlock()

	delete(mm.Msgs, msgID)
}

//定时把没有收到ACK确认的消息插入到数据库中
func (mm *MessageManager) AsyncLoop() {
	for {
		time.Sleep(ASYNC_INTERNAL)
		mm.Sync()
	}
}

func (mm *MessageManager) Sync() {
	mm.lock.Lock()
	defer mm.lock.Unlock()

	for msgID, msg := range mm.Msgs {
		log.Debugf("sync %+v", msg)
		err := model.MessageModel().Insert(msg)
		if nil != err {
			log.Errorln(err) //only log
			continue
		}
		delete(mm.Msgs, msgID)
	}
}
func (mm *MessageManager) SyncByAccount(account string) {
	mm.lock.Lock()
	defer mm.lock.Unlock()

	for msgID, msg := range mm.Msgs {
		if account == msg.Key() {
			log.Debugf("SyncByAccount %+v", msg)
			err := model.MessageModel().Insert(msg)
			if nil != err {
				log.Errorln(err) //only log
				continue
			}
			delete(mm.Msgs, msgID)
		}
	}
}
