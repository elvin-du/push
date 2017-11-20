package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"gokit/util"
	"log"
	"net/http"
)

var E_PUSH_FAILED = errors.New("push failed")

type Client struct {
	RC4Key    string
	AppID     string
	AppSecret string
	UrlStr    string
	bearer    string
}

func NewClient(key, appID, appSecret, urlStr string) (*Client, error) {
	bearer, err := util.RC4EncryptToBase64(key, []byte(appID+":"+appSecret))
	if nil != err {
		log.Println(err)
		return nil, err
	}

	return &Client{
		RC4Key:    key,
		AppID:     appID,
		AppSecret: appSecret,
		UrlStr:    urlStr,
		bearer:    "Bearer " + bearer,
	}, nil
}

func (cli *Client) Push(notification *Notification) error {
	bin, err := json.Marshal(notification)
	if nil != err {
		log.Println(err)
		return err
	}
	log.Println(string(bin))

	r := bytes.NewReader(bin)
	req, err := http.NewRequest("POST", cli.UrlStr, r)
	if nil != err {
		log.Println(err)
		return err
	}

	req.Header.Add("Authorization", cli.bearer)
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if nil != err {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()

	if 200 != resp.StatusCode {
		log.Printf("%+v", resp)
		return E_PUSH_FAILED
	}

	return nil
}
