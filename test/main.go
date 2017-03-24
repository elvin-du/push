package main

import (
	"log"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	err := Connect()
	if nil != err {
		log.Println(err)
		return
	}

	go ReadLoop()

	go Ping()
	select {}
}
