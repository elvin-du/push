package db

import (
	"hscore/log"
	dbredis "push/common/db/redis"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var (
	redis *dbredis.Pool
	mysql *gorm.DB
)

//TODO
func Init() {
	redis = dbredis.NewPool("localhost:6379")
	var err error = nil
	mysql, err = gorm.Open("mysql", "root:JTabc.123@/push_core?charset=utf8&parseTime=True&loc=Local")
	if nil != err {
		log.Fatalln(err)
	}
}

func Redis() *dbredis.Pool {
	return redis
}

func Mysql() *gorm.DB {
	return mysql
}
