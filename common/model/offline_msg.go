package model

import (
	"fmt"
	"gokit/log"
	libdb "push/common/db"
)

type offlineMsg struct{}

var _offlineMsg *offlineMsg

func OfflineMsgModel() *offlineMsg {
	return _offlineMsg
}

func (om *offlineMsg) Get(appID, clientID string) ([]*OfflineMsg, error) {
	key := fmt.Sprintf("%s:%s", appID, clientID)
	db, err := libdb.ShardMysql(key)
	if nil != err {
		log.Errorln(err)
		return nil, err
	}

	sqlStr := fmt.Sprintf("SELECT app_id,client_id,packet_id,kind,content,extra,created_at FROM offline_msgs WHERE app_id='%s' AND client_id='%s' ORDER BY created_at ASC", appID, clientID)
	rows, err := db.Query(sqlStr)
	if nil != err {
		log.Errorln(err, sqlStr)
		return nil, err
	}

	msgs := make([]*OfflineMsg, 0, 0)
	for rows.Next() {
		var msg OfflineMsg
		err = rows.Scan(&msg.AppID, &msg.ClientID, &msg.PacketID, &msg.Kind, &msg.Content, &msg.Extra, &msg.CreateAt)
		if nil != err {
			log.Errorln(err, sqlStr)
			return nil, err
		}

		msgs = append(msgs, &msg)
	}

	return msgs, nil
}
