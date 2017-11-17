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

	err = SingIn("63163c7b40f2abee", "283abdfc9123987980d8aabaa7108e6c", "5a0ea86008b62f0928970a52")
	if nil != err {
		log.Println(err)
		return
	}
	err = SingIn("63163c7b40f2abee", "283abdfc9123987980d8aabaa7108e6c", "5a0ea86508b62f0928970a53")
	if nil != err {
		log.Println(err)
		return
	}

	go ReadLoop()

	go Ping()
	select {}
}
