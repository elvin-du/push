package db

import (
	"database/sql"
	"fmt"
	"gokit/config"
	"gokit/log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"stathat.com/c/consistent"
)

var (
	mysqlDBs = make(map[string]*sql.DB)
	c        = consistent.New()
)

func StartMysql(keys []string) {
	for _, k := range keys {
		startMysql(k)
	}
	log.Infof("mysql db pool:%+v", mysqlDBs)
}

func ShardMysql(key string) (*sql.DB, error) {
	str, err := c.Get(key)
	if nil != err {
		log.Errorln("Shard not found for:  ", key, err)
		return nil, err
	}

	return mysqlDBs[str], nil
}

func MainMysql() (*sql.DB, error) {
	db, ok := mysqlDBs["main"]
	if !ok {
		log.Errorln("Not found main mysqldb")
		return nil, ErrNotFound
	}

	return db, nil
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

	mysqlDBs[key] = NewMysqlDB(addr, useｒName, password, dbName, pool)
	if strings.HasPrefix(key, "shard") {
		c.Add(key)
	}
}

func NewMysqlDB(addr, u, p, dbName string, pool int) *sql.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", u, p, addr, dbName)
	mysqlDB, err := sql.Open("mysql", dsn)
	if nil != err {
		log.Fatalln(err, dsn)
	}
	mysqlDB.SetMaxOpenConns(pool)

	return mysqlDB
}
