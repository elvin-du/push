package main

import (
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

	certBytes := app.Cert
	certPasswd := app.CertPassword
	if isProduction {
		certBytes = app.CertProduction
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

	log.Debugf("%+v %+v %+v %+v\n", res.StatusCode, res.ApnsID, res.Reason, res.Sent())
	return nil
}
