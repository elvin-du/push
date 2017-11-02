package main

import (
	"fmt"
	"gokit/util"

	"stathat.com/c/consistent"
)

func main() {
	//	GetShard("aaaa")
	//	GetAppIDAndAppSecret()
	AESEncryptAppIDAndAppSecret("01e9175ca8805cc2137c44eb86184922")
}

func AESEncryptAppIDAndAppSecret(key string) string {
	id, secret := GetAppIDAndAppSecret()
	data := id + ":" + secret
	fmt.Println("raw data:", data)
	str, err := util.RC4EncryptToBase64(key, []byte(data))
	if nil != err {
		fmt.Println(err)
		return ""
	}

	fmt.Println("RC4 Encrypt:", str)

	bin, err := util.RC4DecryptFromBase64(key, str)
	if nil != err {
		fmt.Println(err)
		return ""
	}

	fmt.Println("RC4 Decrypt:", string(bin))

	return str
}

func GetAppIDAndAppSecret() (id string, secret string) {
	var err error = nil
	id, err = util.RandomString(16)
	if nil != err {
		fmt.Println(err)
		return
	}
	fmt.Println("ID:", id)

	secret, err = util.RandomString(32)
	if nil != err {
		fmt.Println(err)
		return
	}

	fmt.Println("Secret:", secret)
	return id, secret
}

func GetShard(name string) {
	hosts := []string{"shard1", "shard2", "shard3", "shard4", "shard5", "shard6", "shard7", "shard8", "shard9"}
	c := consistent.New()
	for _, h := range hosts {
		c.Add(h)
	}

	h, err := c.Get(name)
	if nil != err {
		fmt.Println(err)
		return
	}

	fmt.Println(h)
	fmt.Printf("SELECT * FROM push_%s.offline_msgs;\n", h)
}
