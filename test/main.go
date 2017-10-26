package main

import (
	"log"

	//	"os"
)

//var f *os.File

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	var err error = nil
	//	f, err = os.Create("log.txt")
	//	if nil != err {
	//		log.Fatal(err)
	//	}

	err = Connect()
	if nil != err {
		log.Println(err)
		return
	}

	go ReadLoop()

	go Ping()
	select {}
}
