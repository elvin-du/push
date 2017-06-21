package session

import (
	redis "push/common/db/redis"
)

type Session struct {
	pool *redis.Pool
}

func New(pool *redis.Pool)*Session{
	return &Session{
		pool:pool,
	}
}

func (s *Session) Get(clientId string) (*session, error) {
	var ses session
	err := s.pool.HMGET(clientId, []interface{}{"client_id", "platform", "gate_server_ip", "gate_server_port"}, &ses)
	if nil != err {
		return nil, err
	}

	return &ses, nil
}
