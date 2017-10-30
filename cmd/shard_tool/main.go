package main

import (
	"fmt"

	"stathat.com/c/consistent"
)

var (
	hosts = []string{"shard1", "shard2", "shard3", "shard4", "shard5", "shard6", "shard7", "shard8", "shard9"}
)

func main() {
	c := consistent.New()
	for _, h := range hosts {
		c.Add(h)
	}

	name := "5001+ANDROID123"
	h, err := c.Get(name)
	if nil != err {
		fmt.Println(err)
		return
	}

	fmt.Println(h)
	fmt.Printf("SELECT * FROM push_%s.offline_msgs;\n", h)
}
