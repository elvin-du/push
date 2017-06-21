package db

import (
	"gokit/log"
	dbredis "push/common/db/redis"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var (
	redis *dbredis.Pool
	mysql *gorm.DB
)

func Start(){
	StartMysql()
	StartRedis()
}

//TODO
func StartRedis() {
	redis = dbredis.NewPool("localhost:6379")
}

func StartMysql(){
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
