package main

import (
	"push/gate/mqtt"
	"sync"
)

type User struct {
	*mqtt.Session
}

func NewUser(ses *mqtt.Session) *User {
	return &User{
		Session: ses,
	}
}

type UserManager struct {
	Users map[string]map[string]*User //app_id:client_id:*User
	Lock  *sync.RWMutex
}

var defaultUserManager = &UserManager{
	Users: make(map[string]map[string]*User),
	Lock:  new(sync.RWMutex),
}

func (um *UserManager) Put(u *User) {
	um.Lock.Lock()
	defer um.Lock.Unlock()

	if _, exist := um.Users[u.Session.AppID]; exist {
		um.Users[u.Session.AppID][u.Session.ClientID] = u
	} else {
		val := map[string]*User{u.Session.ClientID: u}
		um.Users[u.Session.AppID] = val
	}
}

func (um *UserManager) Get(appID, clientID string) *User {
	um.Lock.RLock()
	defer um.Lock.RUnlock()

	if _, exist := um.Users[appID]; exist {
		return um.Users[appID][clientID]
	}

	return nil
}

func (um *UserManager) Remove(appID, clientID string) {
	um.Lock.Lock()
	defer um.Lock.Unlock()

	if _, exist := um.Users[appID]; exist {
		delete(um.Users[appID], clientID)
	}
}
