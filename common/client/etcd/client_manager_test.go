package etcd

import (
	"testing"
)

func TestGetClient(t *testing.T) {
	cli, err := GetClient()
	if nil != err {
		t.Error(err)
		return
	}
	defer Put(cli)

	err = cli.Register("yd", "data", "1.0", "127.0.0.1", "50001", nil)
	if nil != err {
		t.Error(err)
		return
	}

	ip, port, err := cli.Get("yd", "data", "1.0")
	if nil != err {
		t.Error(err)
		return
	}

	t.Log(ip, port)
}

func BenchmarkGetClientParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cli, err := GetClient()
			if nil != err {
				b.Error(err)
				return
			}
			defer Put(cli)

			err = cli.Register("yd", "data", "1.0", "127.0.0.1", "50001", nil)
			if nil != err {
				b.Error(err)
				return
			}

			_, _, err = cli.Get("yd", "data", "1.0")
			if nil != err {
				b.Error(err)
				return
			}
		}
	})
}
