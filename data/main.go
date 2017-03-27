/*
DATA Server提供一切对数据持久化的方法
*/

package main

import (
	"push/gate/service"
)

func main() {
	service.Start()
	StartRPCServer()
}
