package util

import (
	"time"
)

var (
	APPName = "YD"

	EtcdEndpoints     = []string{"127.0.0.1:2379"}
	DialTimeout       = time.Second * 10
	RequestTimeout    = time.Second * 10
	HeartbeatInterval = time.Second * 120 //心跳频率
)
