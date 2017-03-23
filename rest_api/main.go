package main

import (
	"hscore/log"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	StartHTTP()
}
