package util

import (
	"errors"
	"hscore/log"
	"io/ioutil"
	"net"
	"net/http"
)

var (
	InternalIP string
	ExternalIP string
)

func init() {
	var err error = nil
	InternalIP, err = internalIP()
	if nil != err {
		log.Fatalln(err)
	}

	ExternalIP, err = externalIP()
	if nil != err {
		log.Fatalln(err)
	}
}

func internalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if nil != err {
		log.Errorln(err)
		return "", err
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				log.Infoln("local IP:", ipnet.IP.String())
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", errors.New("local ip not found")
}

//重试三次
func externalIP() (string, error) {
	ips := []string{
		"http://ipinfo.io/ip",
		"http://myexternalip.com/raw",
		"http://ipecho.net/plain",
	}

	for _, addr := range ips {
		exIp, err := doExternalIP(addr)
		if nil != err {
			log.Warnln(err)
			continue
		}

		return exIp, nil
	}

	return "", errors.New("Cannot get external IP")
}

func doExternalIP(addr string) (string, error) {
	resp, err := http.Get(addr)
	if nil != err {
		log.Errorln(err)
		return "", err
	}
	defer resp.Body.Close()

	bin, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		log.Errorln(err)
		return "", err
	}

	log.Infoln("external IP:", string(bin))
	return string(bin), nil
}
