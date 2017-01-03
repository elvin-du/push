package main

import (
	"net"
	"testing"
)

func BenchmarkReadWrite(b *testing.B) {
//	for i := 0; i < b.N; i++ {
	for i := 0; i < 10; i++ {
		go func() {
			conn, err := net.Dial("tcp", ":60001")
			if nil != err {
				b.Error(err)
				return
			}
			defer conn.Close()

			_, err = conn.Write([]byte("hi i am tester"))
			if nil != err {
				b.Error(err)
				return
			}

			bin := make([]byte, 1024)
			n, err := conn.Read(bin)
			if nil != err {
				b.Error(err)
				return
			}
			b.Log("read:", string(bin[:n]))
		}()
	}
}
