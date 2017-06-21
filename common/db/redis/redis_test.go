package redis

import (
	"testing"
)

var pool *Pool

func init() {
	pool = NewPool("127.0.0.1:6379")
}

func TestHMSETAndEXPIRE(t *testing.T) {
	m := make(map[string]interface{})
	m["name"] = "elvin"
	m["gender"] = "male"
	m["favorite"] = "htz"
	err := pool.HMSETAndEXPIRE("elvin_user_id", m, 5)
	//	err :=pool.HMSET("elvin_user_id",m)
	if nil != err {
		t.Error(err)
	}
}
