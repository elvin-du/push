package main

import (
	"log"
	"push/gate/service"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	service.Start()
	go defaultServer.StartRPCServer()
	defaultServer.StartTcpServer()
}
