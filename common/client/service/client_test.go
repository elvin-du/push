package service

import (
	"push/meta"
	"testing"

	"golang.org/x/net/context"
)

func TestNewServiceClient(t *testing.T) {
	cli, err := NewServiceClient("116.231.48.149", "50003", "SESSION", "1.0.0")
	if nil != err {
		t.Error(err)
		return
	}

	sessionCli := meta.NewSessionClient(cli.Client)
	infoReq := &meta.SessionInfoRequest{}
	resp,err := sessionCli.Info(context.Background(), infoReq)
	if nil != err {
		t.Error(err)
		return
	}
	t.Log(resp)
}
