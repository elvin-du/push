package util

import (
	"crypto/rc4"
	"encoding/base64"
	"gokit/log"
)

func RC4(key string, data []byte) ([]byte, error) {
	c, err := rc4.NewCipher([]byte(key))
	if nil != err {
		log.Errorln(err)
		return nil, err
	}
	dst := make([]byte, len(data))
	c.XORKeyStream(dst, data)

	return dst, nil
}

func RC4EncryptToBase64(key string, data []byte) (string, error) {
	bin, err := RC4(key, data)
	if nil != err {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(bin), nil
}

func RC4DecryptFromBase64(key, data string) ([]byte, error) {
	dst := make([]byte, len(data))
	n, err := base64.StdEncoding.Decode(dst, []byte(data))
	if nil != err {
		log.Errorln(err)
		return nil, err
	}

	return RC4(key, dst[:n])
}
