package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"flag"
	"gokit/log"
	"gokit/service"
	"gokit/util"
	"io/ioutil"
	"net/http"
	"push/common/model"
)

var (
	devCertFile string
	proCertFile string
	data        string
	method      string
)

func init() {
	flag.StringVar(&devCertFile, "dcf", "", "ios push development cert file path")
	flag.StringVar(&proCertFile, "pcf", "", "ios push production cert file path")
	flag.StringVar(&data, "d", "", "json format data")
	flag.StringVar(&method, "m", "", "method: list create update delete")
}

func main() {
	flag.Parse()
	service.Start()
	InitConf()

	switch method {
	case "list":
	case "create":
		app := &model.App{}
		if "" != data {
			log.Debugf("data:%+v", data)
			err := json.Unmarshal([]byte(data), app)
			if nil != err {
				log.Errorln(err)
				return
			}
		}

		var err error = nil
		if "" == app.ID {
			app.ID, err = util.RandomString(16)
			if nil != err {
				log.Errorln(err)
				return
			}
		}

		if "" == app.Secret {
			app.Secret, err = util.RandomString(32)
			if nil != err {
				log.Errorln(err)
				return
			}
		}

		if "" != devCertFile {
			log.Debugf("devCertFile:%+v", devCertFile)
			bin, err := ioutil.ReadFile(devCertFile)
			if nil != err {
				log.Errorln(err)
				return
			}
			app.Cert = hex.EncodeToString(bin)
		}

		if "" != proCertFile {
			log.Debugf("proCertFile:%+v", proCertFile)
			bin, err := ioutil.ReadFile(proCertFile)
			if nil != err {
				log.Errorln(err)
				return
			}
			app.CertProduction = hex.EncodeToString(bin)
		}

		for _, t := range targets {
			bin, err := json.Marshal(app)
			if nil != err {
				log.Errorln(err)
				return
			}
			r := bytes.NewReader(bin)
			req, err := http.NewRequest("POST", t+"/app", r)
			if nil != err {
				log.Errorln(err)
				return
			}

			resp, err := http.DefaultClient.Do(req)
			if nil != err {
				log.Errorln(err)
				return
			}

			if http.StatusOK != resp.StatusCode {
				log.Errorln("resp:", resp.Status)
				return
			}
		}
		log.Infoln("Create Success")
	case "update":
	case "delete":
	default:
		log.Errorln("unsupport method", method)
	}
}
