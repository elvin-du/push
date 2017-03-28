package dal

import (
	"database/sql"
	"fmt"
	"hscore/log"
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
}

func NewMysql() *Mysql {
	db, err := sql.Open("mysql", "root:@tcp(localhost:4000)/htz_classic")
	if nil != err {
		log.Fatalln(err)
	}

	return &Mysql{DB: db}
}

func (m *Mysql) ReGetDBConn() error {
	db, err := sql.Open("mysql", "root:@tcp(localhost:4000)/htz_classic")
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
	sqlStr := fmt.Sprintf("INSERT INTO %s SET id='%s',gate_server_ip='%s',gate_server_port='%s',user_id='%s',platform='%s',status=1,created_at=%d,updated_at=%d", TBL_CLIENTS, req.ClientId, req.GateIp, req.GatePort, req.UserId, req.Platform, utc, utc)
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
	sqlStr := fmt.Sprintf("UPDATE %s SET status=0,updated_at=%d WHERE client_id='%s'", TBL_CLIENTS, utc, req.ClientId)
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
	sqlStr := fmt.Sprintf("SELECT id,gate_server_ip,gate_server_port,user_id,platform,status,created_at,updated_at FROM %s WHERE user_id='%s'", TBL_CLIENTS, req.UserId)
	log.Debugln(sqlStr)
	rows, err := m.Query(sqlStr)
	if nil != err {
		log.Errorln(err)
		return nil, err
	}

	ret := &meta.GetClientInfoResponse{}
	for rows.Next() {
		item := &meta.GetClientInfoRes{}
		rows.Scan(
			&item.ClientId,
			&item.GateIp,
			&item.GatePort,
			&item.UserId,
			&item.Platform,
			&item.Status,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		ret.UserId = item.UserId
		ret.Items = append(ret.Items, item)
		break //TODO
	}

	return ret, nil
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
	sqlStr := fmt.Sprintf("UPDATE %s SET id='%s',gate_server_ip='%s',gate_server_port='%s',user_id='%s',platform='%s',status=1,updated_at=%d WHERE user_id='%s'", TBL_CLIENTS, req.ClientId, req.GateIp, req.GatePort, req.UserId, req.Platform, utc, req.UserId)
	log.Debugln(sqlStr)
	_, err := m.Query(sqlStr)
	if nil != err {
		log.Errorln(err)
		return nil, err
	}

	return &meta.UpdateClientInfoResponse{}, nil
}
