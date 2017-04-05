package config

import (
	"encoding/json"
	"fmt"
	"hscore/config"
	"hscore/log"
	"hscore/util"
	pushUtil "push/common/util"
)

var (
	RPC_SERVICE_PORT string

	SERVER_IP string

//	MYSQL_DSN  string
//	MYSQL_POOL int
)

func Init() {
	loadConfig()
	ParseConfig()
	StartMysql()
}

func loadConfig() {
	err := config.ReadConfig(util.GetFile("config.yml"))
	if err != nil {
		log.Fatal("Read configuration file failed", err)
	}
}

func ParseConfig() {
	err := config.Get("service:rpc:port", &RPC_SERVICE_PORT)
	if nil != err {
		log.Fatal(err)
	}

	externalIp := false
	err = config.Get("service:externalip", &externalIp)
	if nil != err {
		log.Fatal(err)
	}
	if externalIp {
		SERVER_IP = pushUtil.ExternalIP
	} else {
		SERVER_IP = pushUtil.InternalIP
	}

	//	err = config.Get("db:mysql:dsn", &MYSQL_DSN)
	//	if nil != err {
	//		log.Fatalln(err)
	//	}

	//	err = config.Get("db:mysql:pool", &MYSQL_POOL)
	//	if nil != err {
	//		log.Fatalln(err)
	//	}
}

type YAML_MAP map[interface{}]interface{}
type MysqlSpec struct {
	Addr   string
	User   string
	Passwd string
	DBName string `json:"dbname"`
	Pool   int
}

func (m *MysqlSpec) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4,utf8", m.User, m.Passwd, m.Addr, m.DBName)
}

var (
	MysqlSource = map[string]*MysqlSpec{}
)

func StartMysql() {
	var mysqlSpecs map[interface{}]interface{}
	err := config.Get("mysql", &mysqlSpecs)
	if err != nil {
		log.Fatalln(err)
	}

	log.Debug(mysqlSpecs)
	for name, spec := range mysqlSpecs {
		var ms MysqlSpec
		err := unmarshal(spec.(map[interface{}]interface{}), &ms)
		if nil != err {
			log.Fatalln(err)
		}
		MysqlSource[name.(string)] = &ms
	}
}

func unmarshal(data YAML_MAP, target interface{}) error {
	mp := map[string]interface{}{}
	for k, v := range data {
		mp[k.(string)] = v
	}

	str, err := json.Marshal(mp)
	if err != nil {
		log.Errorln(err)
		return err
	}

	return json.Unmarshal(str, target)
}
