package main

import (
	"errors"
	"gokit/log"
	"gokit/util"
	"push/common/model"
	"push/rest_api/service/config"
	"strings"
)

var (
	E_AUTH_FAILED = errors.New("Auth failed")
)

const (
	Auth_Bearer = "Bearer"
)

func Auth(ctx *Context) error {
	bearer := ctx.Request.Header.Get("Authorization")
	log.Debugln("bearer", bearer)
	if strings.TrimSpace(bearer) == "" {
		log.Errorln("Bearer empty")
		ctx.AbortWithError(401, E_AUTH_FAILED)
		return E_AUTH_FAILED
	}

	if !strings.HasPrefix(bearer, Auth_Bearer) {
		log.Errorln("Auth type invalid")
		ctx.AbortWithError(401, E_AUTH_FAILED)
		return E_AUTH_FAILED
	}

	bearer = strings.TrimSpace(string(bearer[len(Auth_Bearer):]))
	bin, err := util.RC4DecryptFromBase64(config.AUTH_KEY, bearer)
	if nil != err {
		log.Errorln(err)
		ctx.AbortWithError(401, E_AUTH_FAILED)
		return E_AUTH_FAILED
	}
	log.Debugln("Decrypted bearer:", string(bin))

	tmp := strings.Split(string(bin), ":")
	if 2 != len(tmp) {
		log.Errorln("Invalid Bearer:", tmp)
		ctx.AbortWithError(401, E_AUTH_FAILED)
		return E_AUTH_FAILED
	}

	err = model.AuthApp(tmp[0], tmp[1])
	if nil != err {
		log.Errorln(err)
		ctx.AbortWithError(401, E_AUTH_FAILED)
		return E_AUTH_FAILED
	}

	ctx.AppID = tmp[0]
	return nil
}
