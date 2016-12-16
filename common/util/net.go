package util

import (
	"errors"
	"log"
	"net"
)

func LocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if nil != err {
		log.Println(err)
		return "", err
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				log.Println("local ip:", ipnet.IP.String())
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", errors.New("local ip not found")
}
