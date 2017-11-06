package model

import (
	"gokit/log"
	"gokit/util"
	libdb "push/common/db"

	//	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	CERT_PASSWORD_RC4_KEY = "01e9175ca8805cc2137c44eb86184922"
)

type app struct{}

var _app *app

func AppModel() *app {
	return _app
}

func (a *app) GetAll() ([]*App, error) {
	c := libdb.MainMgoDB().C("apps")

	var apps = []*App{}
	it := c.Find(bson.M{}).Iter()
	defer it.Close()

	var app App
	for it.Next(&app) {
		if "" != app.CertPassword {
			bin, err := util.RC4DecryptFromBase64(CERT_PASSWORD_RC4_KEY, app.CertPassword)
			if nil != err {
				log.Errorln(err)
				return nil, err
			}

			app.CertPassword = string(bin)
		}

		if "" != app.CertPasswordProduction {
			bin, err := util.RC4DecryptFromBase64(CERT_PASSWORD_RC4_KEY, app.CertPasswordProduction)
			if nil != err {
				log.Errorln(err)
				return nil, err
			}

			app.CertPasswordProduction = string(bin)
		}

		apps = append(apps, &app)
	}

	return apps, nil
}

func (a *app) Create(app *App) error {
	if "" != app.CertPassword {
		pw, err := util.RC4EncryptToBase64(CERT_PASSWORD_RC4_KEY, []byte(app.CertPassword))
		if nil != err {
			log.Errorln(err)
			return err
		}
		app.CertPassword = pw
	}

	if "" != app.CertPasswordProduction {
		pw, err := util.RC4EncryptToBase64(CERT_PASSWORD_RC4_KEY, []byte(app.CertPasswordProduction))
		if nil != err {
			log.Errorln(err)
			return err
		}
		app.CertPasswordProduction = pw
	}

	c := libdb.MainMgoDB().C("apps")
	err := c.Insert(app)
	if nil != err {
		log.Errorln(err)
		return err
	}

	return nil
}

func (a *app) Update(app *App) error {
	if "" != app.CertPassword {
		pw, err := util.RC4EncryptToBase64(CERT_PASSWORD_RC4_KEY, []byte(app.CertPassword))
		if nil != err {
			log.Errorln(err)
			return err
		}
		app.CertPassword = pw
	}

	if "" != app.CertPasswordProduction {
		pw, err := util.RC4EncryptToBase64(CERT_PASSWORD_RC4_KEY, []byte(app.CertPasswordProduction))
		if nil != err {
			log.Errorln(err)
			return err
		}
		app.CertPasswordProduction = pw
	}

	c := libdb.MainMgoDB().C("apps")
	err := c.UpdateId(app.ID, app)
	if nil != err {
		log.Errorln(err)
		return err
	}

	return nil
}

func (a *app) Delete(appID string) error {
	c := libdb.MainMgoDB().C("apps")
	err := c.RemoveId(appID)
	if nil != err {
		log.Errorln(err)
		return err
	}

	return nil
}
