package etcd

import (
	"testing"
	"time"
)

//func TestETCDRegister(t *testing.T) {
//	_etcd, err := NewETCDClient([]string{"127.0.0.1:2379"}, time.Second*10, time.Second*10)
//	if nil != err {
//		t.Error(err)
//		return
//	}

//	err = _etcd.Register("yd", "data", "1.0", "127.0.0.1", "50001", nil)
//	if nil != err {
//		t.Error(err)
//		return
//	}
//}

func TestETCDGet(t *testing.T) {
	_etcd, err := NewETCDClient([]string{"127.0.0.1:2379"}, time.Second*10, time.Second*10)
	if nil != err {
		t.Error(err)
		return
	}

	ip, port, err := _etcd.Get("GATE", "1.0.0")
	if nil != err {
		t.Error(err)
		return
	}

	t.Log(ip, port)
}
