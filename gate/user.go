package main

import (
	"encoding/json"
	"gokit/log"
	"push/common/model"
	"push/gate/mqtt"
	"push/gate/service/config"
	"push/gate/service/session"

	"github.com/surgemq/message"
)

type User struct {
	*mqtt.Session
	AppID string
	RegID string
}

func NewUser(ses *mqtt.Session, appID, regID string) *User {
	return &User{
		Session: ses,
		AppID:   appID,
		RegID:   regID,
	}
}

func Auth(ses *mqtt.Session, msg *message.ConnectMessage) (*User, error) {
	regID := string(msg.ClientId())
	err := ValidateRegID(regID)
	if nil != err {
		log.Error(err)
		return nil, err
	}

	appID := string(msg.Username())
	log.Debugf("come to connect,app_id: %s,clientid:%s", appID, regID)
	appSecret := string(msg.Password())

	err = doAuth(appID, appSecret)
	if nil != err {
		log.Error(err)
		return nil, err
	}

	u := NewUser(ses, appID, regID)

	//注册到redis上面
	err = u.Online()
	if nil != err {
		log.Error(err)
		return nil, err
	}

	return u, nil
}

func (u *User) Online() error {
	ses2 := &session.Session{
		AppID:          u.AppID,
		RegID:          u.RegID,
		GateServerIP:   config.SERVER_IP,
		GateServerPort: config.RPC_SERVICE_PORT,
	}

	err := session.Update(ses2)
	if nil != err {
		log.Error(err)
		return err
	}

	defaultServer.PutUser(u)
	log.Infof("app_id:%s,client_id:%s online", u.AppID, u.RegID)

	return nil
}

func (u *User) CheckOfflineMsgs() {
	msgs, err := model.OfflineMsgModel().List(u.AppID, u.RegID)
	if nil != err {
		log.Errorln(err)
		return
	}

	log.Debugf("found %d offline msg for app_id:%s,reg_id:%s", len(msgs), u.AppID, u.RegID)
	for _, v := range msgs {
		msg := Message{}
		msg.Content = v.Content
		msg.ID = v.ID
		bin, err := json.Marshal(msg)
		if nil != err {
			log.Errorln(err)
			return
		}

		go u.Push(bin)
	}
}

//TODO auth
func doAuth(appID, regID string) error {
	return nil
}
