package dal

import (
	"database/sql"
	"fmt"
	"hscore/log"
	"push/data/service/config"
	. "push/errors"
	"push/meta"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	TBL_CLIENTS      = "clients"
	TBL_OFFILNE_MSGS = "offline_msgs"
)

type Mysql struct {
	DB *sql.DB
	*config.MysqlSpec
}

func openMysql(spec *config.MysqlSpec) (*sql.DB, error) {
	db, err := sql.Open("mysql", spec.DSN())
	if nil != err {
		log.Errorln(err)
		return nil, err
	}
	db.SetMaxOpenConns(spec.Pool)

	return db, nil
}

func NewMysql(spec *config.MysqlSpec) *Mysql {
	db, err := openMysql(spec)
	if nil != err {
		log.Errorln(err)
	}

	return &Mysql{DB: db, MysqlSpec: spec}
}

func (m *Mysql) ReGetDBConn() error {
	db, err := openMysql(m.MysqlSpec)
	if nil != err {
		log.Errorln(err)
		return err
	}
	m.DB = db

	return nil
}

func (m *Mysql) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if nil != m.DB {
		err := m.ReGetDBConn()
		if nil != err {
			log.Errorln(err)
			return nil, err
		}
	}

	return m.DB.Query(query, args...)
}

var (
	tmpDB = make(map[string]*ClientInfo)
)

//上线
func (m *Mysql) Online(req *meta.DataOnlineRequest) (*meta.DataOnlineResponse, error) {
	utc := time.Now().Unix()
	sqlStr := fmt.Sprintf("INSERT INTO %s SET id='%s',gate_server_ip='%s',gate_server_port='%s',platform='%s',status=1,created_at=%d,updated_at=%d", TBL_CLIENTS, req.ClientId, req.GateServerIP, req.GateServerPort, req.Platform, utc, utc,)
	log.Debugln(sqlStr)

	_, err := m.Query(sqlStr)
	if nil != err {
		log.Errorln(err)
		return nil, err
	}

	return &meta.DataOnlineResponse{}, nil
}

//下线
func (m *Mysql) Offline(req *meta.DataOfflineRequest) (*meta.DataOfflineResponse, error) {
	utc := time.Now().Unix()
	sqlStr := fmt.Sprintf("UPDATE %s SET status=0,updated_at=%d WHERE id='%s'", TBL_CLIENTS, utc, req.ClientId)
	log.Debugln(sqlStr)

	_, err := m.Query(sqlStr)
	if nil != err {
		log.Errorln(err)
		return nil, err
	}

	return &meta.DataOfflineResponse{}, nil
}

func (m *Mysql) SaveOfflineMsg(req *meta.SaveOfflineMsgRequest) (*meta.SaveOfflineMsgResponse, error) {
	return &meta.SaveOfflineMsgResponse{}, nil
}

func (m *Mysql) GetOfflineMsgs(req *meta.GetOfflineMsgsRequest) (*meta.GetOfflineMsgsResponse, error) {
	return &meta.GetOfflineMsgsResponse{}, nil
}

func (m *Mysql) DelOfflineMsgs(req *meta.DelOfflineMsgsRequest) (*meta.DelOfflineMsgsResponse, error) {
	return &meta.DelOfflineMsgsResponse{}, nil
}

func (m *Mysql) GetClientInfo(req *meta.GetClientInfoRequest) (*meta.GetClientInfoResponse, error) {
	sqlStr := fmt.Sprintf("SELECT id,gate_server_ip,gate_server_port,platform,status,created_at,updated_at FROM %s WHERE id='%s'", TBL_CLIENTS, req.ClientId)
	log.Debugln(sqlStr)
	rows, err := m.Query(sqlStr)
	if nil != err {
		log.Errorln(err)
		return nil, err
	}

	ret := &meta.GetClientInfoResponse{}
	for rows.Next() {
		rows.Scan(
			&ret.ClientId,
			&ret.GateServerIP,
			&ret.GateServerPort,
			&ret.Platform,
			&ret.Status,
			&ret.CreatedAt,
			&ret.UpdatedAt,
		)
		//以防万一，只取一个
		return ret, nil
	}

	return nil, CLIENT_NOT_FOUND
}

//TODO SQL注入
//func (m *Mysql)getClientInfoById(cliId string)(*Client,error){
//    	sqlStr := fmt.Sprintf("SELECT id,gate_server_ip,user_id,platform,status,created_at,updated_at FROM %s WHERE user_id='%s'", TBL_CLIENTS, cliId)
//	log.Println(sqlStr)
//	rows, err := m.Query(sqlStr)
//	if nil != err {
//		log.Println(err)
//		return nil, err
//	}

//    return
//}

func (m *Mysql) UpdateClientInfo(req *meta.UpdateClientInfoRequest) (*meta.UpdateClientInfoResponse, error) {
	utc := time.Now().Unix()
	sqlStr := fmt.Sprintf("UPDATE %s SET gate_server_ip='%s',gate_server_port='%s',platform='%s',status=1,updated_at=%d WHERE id='%s'", TBL_CLIENTS, req.GateServerIP, req.GateServerPort, req.Platform, utc, req.ClientId)
	log.Debugln(sqlStr)
	_, err := m.Query(sqlStr)
	if nil != err {
		log.Errorln(err)
		return nil, err
	}

	return &meta.UpdateClientInfoResponse{}, nil
}
