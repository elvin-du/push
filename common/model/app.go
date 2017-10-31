package model

import (
	"fmt"
	"gokit/log"
	libdb "push/common/db"
)

type app struct{}

var _app *app

func AppModel() *app {
	return _app
}

func (a *app) GetAll() ([]*App, error) {
	db, err := libdb.MainMysql()
	if nil != err {
		log.Errorln(err)
		return nil, err
	}
	defer db.Close()

	sqlStr := fmt.Sprintf("SELECT id,secret,auth_type,name,description,status,created_at,updated_at FROM apps WHERE status=1")
	rows, err := db.Query(sqlStr)
	if nil != err {
		log.Errorln(err, sqlStr)
		return nil, err
	}

	apps := make([]*App, 0, 0)
	for rows.Next() {
		var app App
		err = rows.Scan(&app.ID, &app.Secret, &app.AuthType, &app.Name, &app.Description, &app.Status, &app.CreateAt, &app.UpdatedAt)
		if nil != err {
			log.Errorln(err, sqlStr)
			return nil, err
		}

		apps = append(apps, &app)
	}

	return apps, nil
}
