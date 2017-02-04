package main

import (
	"push/gate/service"
)

func main() {
	service.Start()
	go defaultServer.StartRPCServer()
	defaultServer.StartTcpServer()
}
