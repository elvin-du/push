package db

import (
	"database/sql"
	"fmt"
	"gokit/config"
	"gokit/log"
	"push/common/db"

	"strings"

	"stathat.com/c/consistent"
)

var (
	mysqlDBs = make(map[string]*sql.DB)
	c        = consistent.New()
)

func StartMysql() {
	startMysql("main")
	startMysql("shard1")
	startMysql("shard2")
	startMysql("shard3")
	startMysql("shard4")
	startMysql("shard5")
	startMysql("shard6")
	startMysql("shard7")
	startMysql("shard8")
	startMysql("shard9")
}

func ShardMysql(key string) (*sql.DB, error) {
	str, err := c.Get(key)
	if nil != err {
		log.Errorln("Shard not found for:  ", key, err)
		return nil, err
	}

	return mysqlDBs[str], nil
}

func MainMysql() *sql.DB {
	return mysqlDBs["main"]
}

func startMysql(key string) {
	var (
		addr     string
		useｒName string
		password string
		dbName   string
		pool     int
	)

	err := config.Get(fmt.Sprintf("mysql:%s:addr", key), &addr)
	if nil != err {
		log.Fatal(err)
	}
	err = config.Get(fmt.Sprintf("mysql:%s:user", key), &useｒName)
	if nil != err {
		log.Fatal(err)
	}
	err = config.Get(fmt.Sprintf("mysql:%s:passwd", key), &password)
	if nil != err {
		log.Fatal(err)
	}
	err = config.Get(fmt.Sprintf("mysql:%s:dbname", key), &dbName)
	if nil != err {
		log.Fatal(err)
	}
	err = config.Get(fmt.Sprintf("mysql:%s:pool", key), &pool)
	if nil != err {
		log.Fatal(err)
	}

	mysqlDBs[key] = db.NewMysqlDB(addr, useｒName, password, dbName, pool)
	if strings.HasPrefix(key, "shard") {
		c.Add(key)
	}
}
