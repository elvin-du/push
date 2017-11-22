package model

import (
	"gokit/log"
	libdb "push/common/db"
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
	it := c.Find(nil).Iter()
	defer it.Close()

	var app App
	for it.Next(&app) {
		err := app.Decrypt()
		if nil != err {
			log.Errorln(err)
			return nil, err
		}

		apps = append(apps, &app)
	}

	return apps, nil
}

func (a *app) Create(app *App) error {
	err := app.Encrypt()
	if nil != err {
		log.Errorln(err)
		return err
	}

	c := libdb.MainMgoDB().C("apps")
	err = c.Insert(app)
	if nil != err {
		log.Errorln(err)
		return err
	}

	return nil
}

func (a *app) Update(app *App) error {
	err := app.Encrypt()
	if nil != err {
		log.Errorln(err)
		return err
	}

	c := libdb.MainMgoDB().C("apps")
	err = c.UpdateId(app.ID, app)
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

func (a *app) Get(id string) (*App, error) {
	c := libdb.MainMgoDB().C("apps")

	it := c.FindId(id).Iter()
	defer it.Close()

	var app App
	for it.Next(&app) {
		err := app.Decrypt()
		if nil != err {
			log.Errorln(err)
			return nil, err
		}

		return &app, nil
	}

	return nil, E_NOT_FOUND
}
