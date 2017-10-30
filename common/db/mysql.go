package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func NewMysqlDB(addr, u, p, dbName string, pool int) *sql.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", u, p, addr, dbName)
	mysqlDB, err := sql.Open("mysql", dsn)
	if nil != err {
		log.Fatalln(err, dsn)
	}
	mysqlDB.SetMaxOpenConns(pool)

	return mysqlDB
}
