package model

import (
	"fmt"
	"gokit/log"
	"gokit/util"
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
	db, err := libdb.MainMysql()
	if nil != err {
		log.Errorln(err)
		return nil, err
	}
	defer db.Close()

	sqlStr := fmt.Sprintf("SELECT id,secret,auth_type,name,description,status,created_at,bundle_id,cert,cert_passwd, cert_production,cert_passwd_production FROM apps WHERE status=1")
	rows, err := db.Query(sqlStr)
	if nil != err {
		log.Errorln(err, sqlStr)
		return nil, err
	}

	apps := make([]*App, 0, 0)
	for rows.Next() {
		var app App
		err = rows.Scan(&app.ID, &app.Secret, &app.AuthType, &app.Name, &app.Description, &app.Status, &app.CreateAt, &app.BundleID, &app.Cert, &app.CertPassword, &app.CertProduction, &app.CertPasswordProduction)
		if nil != err {
			log.Errorln(err, sqlStr)
			return nil, err
		}

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
