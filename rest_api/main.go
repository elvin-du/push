package main

import (
	"push/rest_api/service"
)

func main() {
	service.Start()
	StartHTTP()
}
