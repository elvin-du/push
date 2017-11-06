package model

import (
	"errors"
	"gokit/log"
	"sync"
)

var (
	_apps []*App = nil
)

const (
	AUTH_TYUP_SECRET = 1
)

var (
	mu = &sync.RWMutex{}
)

var (
	E_NOT_FOUND = errors.New("Not found")
)

func InitAppCache() error {
	apps, err := AppModel().GetAll()
	if nil != err {
		log.Fatalln(err)
	}

	_apps = apps
	log.Infof("Init app cache success. %+v", _apps)
	return nil
}

func ReloadAppCache() error {
	mu.Lock()
	defer mu.Unlock()

	apps, err := AppModel().GetAll()
	if nil != err {
		log.Errorln(err)
		return err
	}

	_apps = apps
	return nil
}

func GetApps() ([]*App, error) {
	mu.RLock()
	defer mu.RUnlock()

	return _apps, nil
}

func AppByID(id string) (*App, error) {
	mu.RLock()
	defer mu.RUnlock()

	apps, err := GetApps()
	if nil != err {
		return nil, err
	}

	for _, v := range apps {
		if v.ID == id {
			return v, nil
		}
	}

	return nil, E_NOT_FOUND
}

func AuthApp(id, secret string) error {
	mu.RLock()
	defer mu.RUnlock()

	apps, err := GetApps()
	if nil != err {
		return err
	}

	for _, v := range apps {
		if v.ID == id && v.Secret == secret && AUTH_TYUP_SECRET == v.AuthType {
			return nil
		}
	}

	return E_NOT_FOUND
}
