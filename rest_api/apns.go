package main

import (
	"encoding/hex"
	"gokit/log"
	"push/common/model"
	"time"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"github.com/sideshow/apns2/payload"
)

const (
	C_IOS_DEFAULT_SOUND = "default"
)

func IOSPush(appID, devToken, alert, sound string, isProduction bool, badge, ttl uint64, extras map[string]interface{}) (string, error) {
	if "" == sound {
		sound = C_IOS_DEFAULT_SOUND
	}

	app, err := model.AppByID(appID)
	if nil != err {
		log.Errorln(err)
		return "", err
	}

	certBytes, err := hex.DecodeString(app.Cert)
	if nil != err {
		log.Errorln(err)
		return "", err
	}

	certPasswd := app.CertPassword
	if isProduction {
		certBytes, err = hex.DecodeString(app.CertProduction)
		if nil != err {
			log.Errorln(err)
			return "", err
		}

		certPasswd = app.CertPasswordProduction
	}
	cert, err := certificate.FromP12Bytes(certBytes, certPasswd)
	if err != nil {
		log.Errorln("Cert Error:", err)
		return "", err
	}

	client := apns2.NewClient(cert).Development()
	if isProduction {
		client = apns2.NewClient(cert).Production()
	}

	notification := &apns2.Notification{}
	notification.DeviceToken = devToken
	notification.Topic = app.BundleID
	notification.Expiration = time.Now().Add(time.Second * time.Duration(ttl))
	p := payload.NewPayload().Alert(alert).Sound(sound).Badge(int(badge))
	for k, v := range extras {
		p.Custom(k, v)
	}
	notification.Payload = p

	res, err := client.Push(notification)
	if err != nil {
		log.Errorln(err)
		return "", err
	}

	log.Infof("app_id:%s,dev_token:%s,is_production:%+v,status_code:%+v,apns_id:%+v,reason:%+v,ok:%+v", appID, devToken, isProduction, res.StatusCode, res.ApnsID, res.Reason, res.Sent())
	return res.ApnsID, nil
}
