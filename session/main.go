package main

import (
	"push/gate/service"
)

func main() {
	service.Start()
	StartRPCServer()
}
