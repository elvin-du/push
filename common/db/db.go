package db

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

type Option struct {
	MaxIdle        int
	MaxActive      int
	IdleTimeout    time.Duration
	ConnectTimeout time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
}

const (
	MAX_IDLE        = 10
	MAX_ACTIVE      = 100
	IDLE_TIMEOUT    = 300 * time.Second
	CONNECT_TIMEOUT = 30 * time.Second
	READ_TIMEOUT    = 30 * time.Second
	WRITE_TIMEOUT   = 30 * time.Second
)

type Pool struct {
	*redis.Pool
}

func NewPool(addr string) *Pool {
	pool := &redis.Pool{
		MaxIdle:     MAX_IDLE,
		MaxActive:   MAX_ACTIVE,
		IdleTimeout: IDLE_TIMEOUT,
		Dial: func() (redis.Conn, error) {
			return redis.DialTimeout("tcp", addr, CONNECT_TIMEOUT, READ_TIMEOUT, WRITE_TIMEOUT)
		},
	}

	return &Pool{pool}
}

func NewPoolWithOpt(addr string, opt *Option) *Pool {
	pool := &redis.Pool{
		MaxIdle:     opt.MaxIdle,
		MaxActive:   opt.MaxActive,
		IdleTimeout: opt.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			return redis.DialTimeout("tcp", addr, opt.ConnectTimeout, opt.ReadTimeout, opt.WriteTimeout)
		},
	}

	return &Pool{pool}
}

func (p *Pool) EXPIRE(key string, TTL time.Duration) error {
	_, err := p.Get().Do("EXPIRE", key, TTL)
	if nil != err {
		return err
	}

	return nil
}

func (p *Pool) HMSET(key string, fields map[string]interface{}) error {
	args := []interface{}{key}
	for k, v := range fields {
		args = append(args, k)
		args = append(args, v)
	}
	_, err := p.Get().Do("HMSET", args...)
	if nil != err {
		return err
	}

	return nil
}

func (p *Pool) HMSETAndEXPIRE(key string, fields map[string]interface{}, TTL time.Duration) error {
	args := []interface{}{key}
	for k, v := range fields {
		args = append(args, k)
		args = append(args, v)
	}

	err := p.Get().Send("HMSET", args...)
	if nil != err {
		return err
	}

	_, err = p.Get().Do("EXPIRE", key, TTL)
	if nil != err {
		return err
	}

	return err
}

/*
v:must be ptr of struct. and should use redis tag for struct field.
for example:
{
	"name":"elvin" `"redis":"name"`
}
*/
func (p *Pool) HGETALL(key string, v interface{}) error {
	valueInterfaces, err := redis.Values(p.Get().Do("HGETALL", key))
	if nil != err {
		return err
	}

	return redis.ScanStruct(valueInterfaces, &v)
}

/*
v:must be ptr of struct. and should use redis tag for struct field.
for example:
{
	"name":"elvin" `"redis":"name"`
}
*/
func (p *Pool) HMGET(key string, fields []interface{}, v interface{}) error {
	tmp := []interface{}{key}
	tmp = append(tmp, fields...)
	valueInterfaces, err := redis.Values(p.Get().Do("HMGET", tmp...))
	if nil != err {
		return err
	}

	return redis.ScanStruct(valueInterfaces, &v)
}

func (p *Pool) DEL(keys []interface{}) error {
	number, err := redis.Int(p.Get().Do("DEL", keys...))
	if nil != err {
		return err
	}
	if len(keys) != number {
		return fmt.Errorf("deleted %d keys,but expected delete %d keys", number, len(keys))
	}

	return nil
}
