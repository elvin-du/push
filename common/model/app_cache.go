package model

import (
	"errors"
	"sync"
)

var (
	_apps []*App = nil
)

var (
	mu = &sync.RWMutex{}
)

var (
	E_NOT_FOUND = errors.New("Not found")
)

func LoadAppCache() error {
	mu.Lock()
	defer mu.Unlock()

	apps, err := AppModel().GetAll()
	if nil != err {
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
		if v.ID == id && v.Secret == secret {
			return nil
		}
	}

	return E_NOT_FOUND
}
