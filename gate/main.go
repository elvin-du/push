package main

func main() {
	go StartRPCServer()
	StartTcpServer()
}
