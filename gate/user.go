package main

import (
	"push/gate/mqtt"
	"sync"
)

type User struct {
	Session *mqtt.Session
}

func (u *User) NewUser(ses *mqtt.Session) *User {
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

	um.Users[u.AppID][u.ClientId] = u
}

func (um *UserManager) Get(appID, clientID string) *User {
	um.Lock.RLock()
	defer um.Lock.RUnlock()

	return um.Users[u.AppID][u.ClientId]
}

func (um *UserManager) Remove(appID, clientID string) {
	um.Lock.Lock()
	defer um.Lock.Unlock()

	delete(um.Users[u.AppID], clientID)
}
