package main

import (
	"errors"
)

var (
	APP_ID_INVALID       = errors.New("app_id invalid")
	REG_ID_INVALID       = errors.New("reg_id invalid")
	DEV_TOKEN_INVALID    = errors.New("dev_token invalid")
	COTENT_INVALID       = errors.New("content invalid")
	PLATFORM_INVALID     = errors.New("platform invalid")
	REQUEST_DATA_INVALID = errors.New("request data invalid")
)
