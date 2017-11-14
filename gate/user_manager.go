package main

import (
	"sync"
)

type UserManager struct {
	Users map[string]map[string]*User //app_id:reg_id:*User
	//	SessionUsers map[string][]*User          //key: session ID
	Lock *sync.RWMutex
}

func NewUserManager() *UserManager {
	return &UserManager{
		Users: make(map[string]map[string]*User),
		Lock:  new(sync.RWMutex),
	}
}

func (um *UserManager) Put(u *User) {
	um.Lock.Lock()
	defer um.Lock.Unlock()

	if _, exist := um.Users[u.AppID]; exist {
		um.Users[u.AppID][u.RegID] = u
	} else {
		val := map[string]*User{u.RegID: u}
		um.Users[u.AppID] = val
	}
}

func (um *UserManager) Get(appID, regID string) *User {
	um.Lock.RLock()
	defer um.Lock.RUnlock()

	if _, exist := um.Users[appID]; exist {
		return um.Users[appID][regID]
	}

	return nil
}

func (um *UserManager) GetByID(sessionID string) []*User {
	um.Lock.RLock()
	defer um.Lock.RUnlock()

	users := []*User{}
	for _, appUsers := range um.Users {
		for _, user := range appUsers {
			if sessionID == user.ID {
				users = append(users, user)
			}
		}
	}

	if 0 != len(users) {
		return users
	}

	return nil
}

func (um *UserManager) Remove(appID, regID string) *User {
	um.Lock.Lock()
	defer um.Lock.Unlock()

	if _, exist := um.Users[appID]; exist {
		delete(um.Users[appID], regID)
		return um.Users[appID][regID]
	}

	return nil
}

func (um *UserManager) RemoveByID(sessionID string) *User {
	um.Lock.Lock()
	defer um.Lock.Unlock()

	for _, appUsers := range um.Users {
		for _, user := range appUsers {
			if sessionID == user.ID {
				delete(um.Users[user.AppID], user.RegID)
				return user
			}
		}
	}

	return nil
}
