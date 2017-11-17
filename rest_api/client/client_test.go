package client

import (
	"testing"
)

func TestPush(t *testing.T) {
	key := "01e9175ca8805cc2137c44eb86184922"
	appID := "63163c7b40f2abee"
	appSecret := "283abdfc9123987980d8aabaa7108e6c"
	urlStr := "http://localhost:52001/push"
	cli, err := NewClient(key, appID, appSecret, urlStr)
	if nil != err {
		t.Error(err)
		return
	}

	n := &Notification{}
	n.Alert = "hi test"
	n.Audience = "5a0ea86008b62f0928970a52"
	n.TTL = 60 * 60 * 24
	n.AddExtra("kind", 2)
	err = cli.Push(n)
	if nil != err {
		t.Error(err)
	}
}
