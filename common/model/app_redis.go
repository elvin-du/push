package model

import (
	"errors"
	"gokit/log"
	"push/common/db"

	"gopkg.in/mgo.v2/bson"
)

const (
	APP_REDIS_HASH_KEY = "push.redis.apps"
)

var (
	E_NOT_FOUND = errors.New("Not found")
)

func InitAppCache() {
	apps, err := AppModel().GetAll()
	if nil != err {
		log.Warnf("init app cache failed,err:%s", err.Error())
		return
	}

	err = SetAllApps2Cache(apps)
	if nil != err {
		log.Warnf("init app cache failed,err:%s", err.Error())
		return
	}

	log.Infof("Init app cache success")
}

func AppByIDFromCache(id string) (*App, error) {
	bin, err := db.MainRedis().HGET(APP_REDIS_HASH_KEY, id)
	if nil != err {
		log.Warnf("err:%s,id:%s", err.Error(), id)
		return nil, err
	}

	var app App
	err = bson.Unmarshal(bin, &app)
	if nil != err {
		log.Warnf("err:%+v,bin:%+v", err, string(bin))
		return nil, err
	}

	err = app.Decrypt()
	if nil != err {
		log.Warnln(err)
		return nil, err
	}

	return &app, nil
}

func AppByID(id string) (*App, error) {
	app, err := AppByIDFromCache(id)
	if nil == err {
		return app, nil
	}

	app, err = AppModel().Get(id)
	if nil != err {
		log.Errorln(err, "id:", id)
		return nil, err
	}

	SetApp2Cache(app)

	return app, nil
}

func SetApp2Cache(app *App) error {
	err := app.Encrypt()
	if nil != err {
		log.Warnln(err)
		return err
	}

	bin, err := bson.Marshal(app)
	if nil != err {
		log.Warnf("err:%+v,app:%+v", err, app)
		return err
	}

	data := make(map[string]interface{})
	data[app.ID] = bin
	return db.MainRedis().HMSET(APP_REDIS_HASH_KEY, data)
}

func SetApp(app *App) error {
	err := AppModel().Create(app)
	if nil != err {
		log.Errorln(err)
		return err
	}

	SetApp2Cache(app)
	return nil
}

func GetAppsFromCache() ([]*App, error) {
	arrayBytes, err := db.MainRedis().HGETALL2Bytes(APP_REDIS_HASH_KEY)
	if nil != err {
		log.Warnln(err)
		return nil, err
	}

	size := len(arrayBytes) / 2
	apps := make([]*App, size)
	for i := 0; i < size; i++ {
		var a App
		err = bson.Unmarshal(arrayBytes[(i*2)+1], &a)
		if nil != err {
			log.Warnln(err)
			return nil, err
		}

		err = a.Decrypt()
		if nil != err {
			log.Warnln(err)
			return nil, err
		}

		apps = append(apps, &a)
	}

	return apps, nil
}

func SetAllApps2Cache(apps []*App) error {
	data := make(map[string]interface{}, len(apps))
	for _, app := range apps {
		err := app.Encrypt()
		if nil != err {
			log.Warnln(err)
			return err
		}

		bin, err := bson.Marshal(app)
		if nil != err {
			log.Warnf("err:%+v,app:%+v", err, app)
			return err
		}

		data[app.ID] = bin
	}

	return db.MainRedis().HMSET(APP_REDIS_HASH_KEY, data)
}

func GetApps() ([]*App, error) {
	apps, err := GetAppsFromCache()
	if nil == err {
		return apps, nil
	}

	apps, err = AppModel().GetAll()
	if nil != err {
		log.Errorln(err)
		return nil, err
	}

	SetAllApps2Cache(apps)

	return apps, nil
}

func DeleteAppFromCache(id string) error {
	return db.MainRedis().DEL([]interface{}{id})
}

func AuthApp(id, secret string) error {
	app, err := AppByID(id)
	if nil != err {
		log.Errorf("err:%+v,id:%+v", err, id)
		return err
	}

	if id == app.ID && secret == app.Secret {
		return nil
	}

	return E_NOT_FOUND
}
