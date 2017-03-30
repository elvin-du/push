package main

import (
	"push/session/service"
)

func main() {
	service.Start()
	StartRPCServer()
}
