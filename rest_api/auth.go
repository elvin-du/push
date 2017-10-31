package main

import (
	"errors"
	"gokit/log"
	"gokit/util"
	"push/rest_api/service/config"
	"strings"

	"push/common/model"

	"github.com/gin-gonic/gin"
)

var (
	E_AUTH_FAILED = errors.New("Auth failed")
)

func Auth(ctx *gin.Context) error {
	bearer := ctx.Request.Header.Get("Bearer")
	log.Debugln("bearer", bearer)
	if strings.TrimSpace(bearer) == "" {
		log.Errorln("Bearer empty")
		return E_AUTH_FAILED
	}

	bin, err := util.AesDecryptFromHex([]byte(config.AUTH_KEY), bearer)
	if nil != err {
		log.Errorln(err)
		return E_AUTH_FAILED
	}
	log.Debugln("Decrypted bearer:", string(bin))

	tmp := strings.Split(string(bin), ":")
	if 2 != len(tmp) {
		log.Errorln("Invalid Bearer:", tmp)
		return E_AUTH_FAILED
	}

	AuthApp(tmp[0], tmp[1])
	return nil
}

func AuthApp(id, secret string) error {
	apps, err := model.GetApps()
	if nil != err {
		log.Errorln(err)
		return err
	}

	for _, v := range apps {
		if v.ID == id && v.Secret == secret {
			return nil
		}
	}

	return E_AUTH_FAILED
}
