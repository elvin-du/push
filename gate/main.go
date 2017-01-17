package main

import (
	"log"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	go defaultServer.StartRPCServer()
	defaultServer.StartTcpServer()
}
