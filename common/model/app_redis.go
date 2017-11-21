package model

import (
	"errors"
	"gokit/log"
	"push/common/db"

	libRedis "github.com/garyburd/redigo/redis"
	"gopkg.in/mgo.v2/bson"
)

const (
	APP_REDIS_HASH_KEY = "push.redis.apps"
)

var (
	E_NOT_FOUND = errors.New("Not found")
)

func AppByID(id string) (*App, error) {
	bin, err := db.MainRedis().HGET(APP_REDIS_HASH_KEY, id)
	if nil != err && libRedis.ErrNil != err {
		log.Errorln(err, "id:", id)
		return nil, err
	}

	if libRedis.ErrNil == err {
		//		AppModel().
	}

	var app App
	err = bson.UnmarshalJSON(bin, &app)
	if nil != err {
		log.Errorf("err:%+v,bin:%+v", err, string(bin))
		return nil, err
	}

	return &app, nil
}

func SetApp(app *App) error {
	bin, err := bson.Marshal(app)
	if nil != err {
		log.Errorf("err:%+v,app:%+v", err, app)
		return err
	}

	data := make(map[string]interface{})
	data[app.ID] = bin
	return db.MainRedis().HMSET(APP_REDIS_HASH_KEY, data)
}

func GetApps() ([]*App, error) {
	arrayBytes, err := db.MainRedis().HGETALL2Bytes(APP_REDIS_HASH_KEY)
	if nil != err {
		log.Errorln(err)
		return nil, err
	}

	size := len(arrayBytes) / 2
	apps := make([]*App, size)
	for i := 0; i < size; i++ {
		var a App
		err = bson.UnmarshalJSON(arrayBytes[(i*2)+1], &a)
		if nil != err {
			log.Errorln(err)
			return nil, err
		}
		apps = append(apps, &a)
	}

	return apps, nil
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
