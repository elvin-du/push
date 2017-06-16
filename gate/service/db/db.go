package db

import (
	dbkit "push/common/db"
)

var (
	redis *dbkit.Pool
)

//TODO
func Init() {
	redis = dbkit.NewPool("localhost:6379")
}

func Redis() *dbkit.Pool {
	return redis
}
