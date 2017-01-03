package main

import (
	"log"
	"net"
)

var (
	C_TCP_PORT = ":60001"
)

func main() {
	conn, err := net.Dial("tcp", C_TCP_PORT)
	if nil != err {
		log.Println(err)
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte("hi i am tester"))
	if nil != err {
		log.Println(err)
		return
	}

	log.Println("send done")
	bin := make([]byte, 1024)

	n, err := conn.Read(bin)
	if nil != err {
		log.Println(err)
		return
	}
	log.Println("read:", string(bin[:n]))
}
