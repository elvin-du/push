// +build mysql

package model

import (
	"fmt"
	"gokit/log"
	"gokit/util"
	libdb "push/common/db"
)

type offlineMsg struct{}

var _offlineMsg *offlineMsg

func OfflineMsgModel() *offlineMsg {
	return _offlineMsg
}

func (om *offlineMsg) List(appID, regID string) ([]*OfflineMsg, error) {
	key := fmt.Sprintf("%s:%s", appID, regID)
	db, err := libdb.ShardMysql(key)
	if nil != err {
		log.Errorln(err)
		return nil, err
	}

	sqlStr := fmt.Sprintf("SELECT id,app_id,reg_id,kind,content,extra,created_at FROM offline_msgs WHERE app_id='%s' AND client_id='%s' ORDER BY created_at ASC", appID, clientID)
	rows, err := db.Query(sqlStr)
	if nil != err {
		log.Errorln(err, sqlStr)
		return nil, err
	}

	msgs := make([]*OfflineMsg, 0, 0)
	for rows.Next() {
		var msg OfflineMsg
		err = rows.Scan(&msg.ID, &msg.AppID, &msg.RegID, &msg.Kind, &msg.Content, &msg.Extra, &msg.CreateAt)
		if nil != err {
			log.Errorln(err, sqlStr)
			return nil, err
		}

		msgs = append(msgs, &msg)
	}

	return msgs, nil
}

func (om *offlineMsg) Insert(msg *OfflineMsg) error {
	key := fmt.Sprintf("%s:%s", msg.AppID, msg.RegID)
	db, err := libdb.ShardMysql(key)
	if nil != err {
		log.Errorln(err)
		return err
	}

	sqlStr := fmt.Sprintf("INSERT offline_msgs SET id='%s', app_id='%s',reg_id='%s',kind=%d,content='%s',extra='%s',created_at=%d",
		msg.ID, msg.AppID, msg.RegID, msg.Kind, msg.Content, msg.Extra, util.Timestamp())
	_, err = db.Query(sqlStr)
	if nil != err {
		log.Errorln(err, sqlStr)
		return err
	}

	return nil
}
