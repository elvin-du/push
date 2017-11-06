package db

import (
	"fmt"
	"gokit/config"
	"gokit/log"

	"gopkg.in/mgo.v2"
	//	"gopkg.in/mgo.v2/bson"
)

var (
	mongoSessions = make(map[string]*mgo.Session)
)

func StartMongo(keys []string) {
	for _, k := range keys {
		startMongo(k)
	}
	log.Infof("mongo pool:%+v", mongoSessions)
}

func startMongo(key string) {
	var (
		url string
	)

	err := config.Get(fmt.Sprintf("mongo:%s:url", key), &url)
	if nil != err {
		log.Fatal(err)
	}

	ses, err := mgo.Dial(url)
	if nil != err {
		log.Fatal(err)
	}
	mongoSessions[key] = ses
}

func MainMgo() *mgo.Session {
	return mongoSessions["main"]
}

func MainMgoDB() *mgo.Database {
	return MainMgo().DB("push_core")
}
