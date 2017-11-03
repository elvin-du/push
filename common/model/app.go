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

	sqlStr := fmt.Sprintf("SELECT id,secret,auth_type,name,description,status,created_at,updated_at,bundle_id,cert,cert_passwd, cert_production,cert_passwd_production FROM apps WHERE status=1")
	rows, err := db.Query(sqlStr)
	if nil != err {
		log.Errorln(err, sqlStr)
		return nil, err
	}

	apps := make([]*App, 0, 0)
	for rows.Next() {
		var app App
		err = rows.Scan(&app.ID, &app.Secret, &app.AuthType, &app.Name, &app.Description, &app.Status, &app.CreatedAt, &app.UpdatedAt, &app.BundleID, &app.Cert, &app.CertPassword, &app.CertProduction, &app.CertPasswordProduction)
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

func (a *app) Create(app *App) error {
	db, err := libdb.MainMysql()
	if nil != err {
		log.Errorln(err)
		return err
	}

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

	sqlStr := fmt.Sprintf(`INSERT apps SET 
		id='%s',secret='%s',auth_type=%d,name='%s',description='%s',
		status=%d,created_at=%d,bundle_id='%s',cert='%s',cert_passwd='%s',
		cert_production='%s',cert_passwd_production='%s',updated_at=%d`,
		app.ID, app.Secret, app.AuthType, app.Name, app.Description,
		app.Status, util.Timestamp(), app.BundleID, app.Cert, app.CertPassword,
		app.CertProduction, app.CertPasswordProduction, 0)
	_, err = db.Query(sqlStr)
	if nil != err {
		log.Errorln(err, sqlStr)
		return err
	}

	return nil
}

func (a *app) Update(app *App) error {
	db, err := libdb.MainMysql()
	if nil != err {
		log.Errorln(err)
		return err
	}

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

	sqlStr := fmt.Sprintf(`UPDATE apps SET 
		id='%s',secret='%s',auth_type=%d,name='%s',description='%s',
		status=%d,updated_at=%d,bundle_id='%s',cert='%s',cert_passwd='%s',
		cert_production='%s',cert_passwd_production='%s'`,
		app.ID, app.Secret, app.AuthType, app.Name, app.Description,
		app.Status, util.Timestamp(), app.BundleID, app.Cert, app.CertPassword,
		app.CertProduction, app.CertPasswordProduction)
	_, err = db.Query(sqlStr)
	if nil != err {
		log.Errorln(err, sqlStr)
		return err
	}

	return nil
}

func (a *app) Delete(appID string) error {
	db, err := libdb.MainMysql()
	if nil != err {
		log.Errorln(err)
		return err
	}

	sqlStr := fmt.Sprintf(`DELETE apps WHERE id='%s'`, appID)
	_, err = db.Query(sqlStr)
	if nil != err {
		log.Errorln(err, sqlStr)
		return err
	}

	return nil
}
