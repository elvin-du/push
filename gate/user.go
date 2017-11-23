package main

import (
	"encoding/json"
	"fmt"
	"gokit/log"
	"gokit/util"
	"push/common/model"
	. "push/errors"
	"push/gate/mqtt"
	"push/gate/service/config"
	"push/gate/service/session"
	"strings"

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
		log.Error(err, "reg_id", regID)
		return nil, err
	}

	appID := string(msg.Username())
	log.Debugf("come to connect,app_id: %s,reg_id:%s", appID, regID)
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
	log.Infof("user app_id:%s,reg_id:%s online", u.AppID, u.RegID)

	return nil
}

func (u *User) CheckOfflineMsgs() {
	msgs, err := model.MessageModel().List(u.AppID, u.RegID)
	if nil != err {
		log.Errorln(err)
		return
	}

	log.Debugf("found %d offline msg for app_id:%s,reg_id:%s", len(msgs), u.AppID, u.RegID)
	for _, v := range msgs {
		msg := PublishMessage{}
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

func (u *User) Key() string {
	return fmt.Sprintf("%s:%s", u.AppID, u.RegID)
}

func doAuth(appID, appSecret string) error {
	bin, err := util.RC4DecryptFromBase64(config.AUTH_KEY, appSecret)
	if nil != err {
		log.Errorln(err)
		return err
	}
	log.Debugln("Decrypted app_secret:", string(bin))

	tmp := strings.Split(string(bin), ":")
	if 2 != len(tmp) {
		log.Errorln("Invalid app_secret:", tmp)
		return REQ_DATA_INVALID
	}

	if appID != tmp[0] {
		log.Errorf("Decrypted app_id:%s != parameter app_id:%s", tmp[0], appID)
		return REQ_DATA_INVALID
	}

	err = model.AuthApp(appID, tmp[1])
	if nil != err {
		log.Errorln(err)
		return err
	}

	return nil
}
