package redis

import (
	"fmt"
	"time"

	libRedis "github.com/garyburd/redigo/redis"
)

var ErrNotFound = fmt.Errorf("Not Found\n")

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
	*libRedis.Pool
}

func NewPool(addr string) *Pool {
	pool := &libRedis.Pool{
		MaxIdle:     MAX_IDLE,
		MaxActive:   MAX_ACTIVE,
		IdleTimeout: IDLE_TIMEOUT,
		Dial: func() (libRedis.Conn, error) {
			return libRedis.DialTimeout("tcp", addr, CONNECT_TIMEOUT, READ_TIMEOUT, WRITE_TIMEOUT)
		},
	}

	return &Pool{pool}
}

func NewPoolWithOpt(addr string, opt *Option) *Pool {
	pool := &libRedis.Pool{
		MaxIdle:     opt.MaxIdle,
		MaxActive:   opt.MaxActive,
		IdleTimeout: opt.IdleTimeout,
		Dial: func() (libRedis.Conn, error) {
			return libRedis.DialTimeout("tcp", addr, opt.ConnectTimeout, opt.ReadTimeout, opt.WriteTimeout)
		},
	}

	return &Pool{pool}
}

func (p *Pool) EXPIRE(key string, TTL int) error {
	c := p.Get()
	defer p.Close()

	_, err := c.Do("EXPIRE", key, TTL)
	if nil != err {
		return err
	}

	return nil
}

func (p *Pool) HMSET(key string, fields map[string]interface{}) error {
	c := p.Get()
	defer p.Close()

	args := []interface{}{key}
	for k, v := range fields {
		args = append(args, k)
		args = append(args, v)
	}
	_, err := c.Do("HMSET", args...)
	if nil != err {
		return err
	}

	return nil
}

func (p *Pool) HMSETAndEXPIRE(key string, fields map[string]interface{}, TTL int) error {
	args := []interface{}{key}
	for k, v := range fields {
		args = append(args, k)
		args = append(args, v)
	}

	c := p.Get()
	defer p.Close()

	err := c.Send("MULTI")
	if nil != err {
		return err
	}

	err = c.Send("HMSET", args...)
	if nil != err {
		return err
	}

	err = c.Send("EXPIRE", key, int(TTL))
	if nil != err {
		return err
	}

	_, err = c.Do("EXEC")
	if nil != err {
		return err
	}

	return nil
}

/*
v:must be ptr of struct. and should use libRedis tag for struct field.
for example:
{
	"name":"elvin" `"redis":"name"`
}
*/
func (p *Pool) HGETALL(key string, v interface{}) error {
	c := p.Get()
	defer p.Close()

	valueInterfaces, err := libRedis.Values(c.Do("HGETALL", key))
	if nil != err {
		return err
	}

	return libRedis.ScanStruct(valueInterfaces, &v)
}

/*
v:must be ptr of struct. and should use redis tag for struct field.
for example:
{
	"name":"elvin" `"redis":"name"`
}
*/
func (p *Pool) HMGET(key string, fields []interface{}, v interface{}) error {
	c := p.Get()
	defer p.Close()

	tmp := []interface{}{key}
	tmp = append(tmp, fields...)
	valueInterfaces, err := libRedis.Values(c.Do("HMGET", tmp...))
	if nil != err {
		return err
	}

	//HMGET在找不到数据的时候，也不会返回err，而是数据返回nil
	for _, v := range valueInterfaces {
		if nil == v {
			return ErrNotFound
		}
	}

	return libRedis.ScanStruct(valueInterfaces, v)
}

func (p *Pool) DEL(keys []interface{}) error {
	c := p.Get()
	defer p.Close()

	number, err := libRedis.Int(c.Do("DEL", keys...))
	if nil != err {
		return err
	}
	if len(keys) != number {
		return fmt.Errorf("deleted %d keys,but expected delete %d keys", number, len(keys))
	}

	return nil
}
