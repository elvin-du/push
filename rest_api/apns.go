package main

import (
	"encoding/hex"
	"gokit/log"
	"push/common/model"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"github.com/sideshow/apns2/payload"
)

func IOSPush(appID, devToken, alert string, isProduction bool) error {
	app, err := model.AppByID(appID)
	if nil != err {
		log.Errorln(err)
		return err
	}

	certBytes, err := hex.DecodeString(app.Cert)
	if nil != err {
		log.Errorln(err)
		return err
	}

	certPasswd := app.CertPassword
	if isProduction {
		certBytes, err = hex.DecodeString(app.CertProduction)
		if nil != err {
			log.Errorln(err)
			return err
		}

		certPasswd = app.CertPasswordProduction
	}
	cert, err := certificate.FromP12Bytes(certBytes, certPasswd)
	if err != nil {
		log.Errorln("Cert Error:", err)
		return err
	}

	client := apns2.NewClient(cert).Development()
	if isProduction {
		client = apns2.NewClient(cert).Production()
	}

	notification := &apns2.Notification{}
	notification.DeviceToken = devToken
	notification.Topic = app.BundleID
	notification.Payload = payload.NewPayload().Alert(alert)

	res, err := client.Push(notification)
	if err != nil {
		log.Errorln(err)
		return err
	}

	log.Infof("app_id:%s,dev_token:%s,is_production:%+v,status_code:%+v,apns_id:%+v,reason:%+v,ok:%+v", appID, devToken, isProduction, res.StatusCode, res.ApnsID, res.Reason, res.Sent())
	return nil
}
