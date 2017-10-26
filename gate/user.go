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

	um.Users[u.Session.AppID][u.Session.ClientID] = u
}

func (um *UserManager) Get(appID, clientID string) *User {
	um.Lock.RLock()
	defer um.Lock.RUnlock()

	return um.Users[appID][clientID]
}

func (um *UserManager) Remove(appID, clientID string) {
	um.Lock.Lock()
	defer um.Lock.Unlock()

	delete(um.Users[appID], clientID)
}
