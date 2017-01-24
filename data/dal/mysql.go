package dal

import (
	"database/sql"
	"fmt"
	"log"
	"push/meta"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

var (
	TBL_CLIENTS      = "clients"
	TBL_OFFILNE_MSGS = "offline_msgs"
)

type Mysql struct {
	DB *sql.DB
}

func NewMysql() *Mysql {
	db, err := sql.Open("mysql", "root:@tcp(localhost:4000)/push_core")
	if nil != err {
		log.Println(err)
	}

	return &Mysql{DB: db}
}

func (m *Mysql) ReGetDBConn() error {
	db, err := sql.Open("mysql", "root:@tcp(localhost:4000)/push_core")
	if nil != err {
		log.Println(err)
		return err
	}
	m.DB = db

	return nil
}

func (m *Mysql) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if nil != m.DB {
		err := m.ReGetDBConn()
		if nil != err {
			log.Println(err)
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
	//TODO
	log.Printf("%s online sucess", req.ClientId)
	//	sqlStr := fmt.Sprintf("INSERT INTO %s SET id='%s',gate_server_ip='%s',user_id='%s',platform='%s',status=1,created_at=%d,updated_at=%d", TBL_CLIENTS, req.ClientId, req.IP, req.UserId, req.Platform, utc, utc)
	//	log.Println(sqlStr)
	//	_, err := m.Query(sqlStr)
	//	if nil != err {
	//		log.Println(err)
	//		return nil, err
	//	}
	ci := &ClientInfo{}
	ci.GateServerIp = req.IP
	ci.Id = req.ClientId
	ci.Platform = req.Platform
	ci.UserId = req.UserId
	ci.Status = 1
	ci.CreatedAt = uint64(utc)
	ci.UpdatedAt = uint64(utc)
	tmpDB[ci.UserId] = ci
	return &meta.DataOnlineResponse{}, nil
}

//下线
func (m *Mysql) Offline(req *meta.DataOfflineRequest) (*meta.DataOfflineResponse, error) {
	utc := time.Now().Unix()
	sqlStr := fmt.Sprintf("UPDATE %s SET status=0,updated_at=%d WHERE client_id='%s'", TBL_CLIENTS, utc, req.ClientId)
	log.Println(sqlStr)

	_, err := m.Query(sqlStr)
	if nil != err {
		log.Println(err)
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
	//	sqlStr := fmt.Sprintf("SELECT id,gate_server_ip,user_id,platform,status,created_at,updated_at FROM %s WHERE user_id='%s'", TBL_CLIENTS, req.UserId,req.)
	//	log.Println(sqlStr)
	//	_, err := m.Query(sqlStr)
	//	if nil != err {
	//		log.Println(err)
	//		return nil, err
	//	}

	ret := &meta.GetClientInfoResponse{}
	v, ok := tmpDB[req.UserId]
	if ok {
		item := &meta.GetClientInfoRes{}
		item.ClientId = v.Id
		item.Platform = v.Platform
		item.IP = v.GateServerIp
		ret.Items = append(ret.Items, item)
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
	return &meta.UpdateClientInfoResponse{}, nil
}
