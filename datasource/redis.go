package datasource

import (
	"github.com/garyburd/redigo/redis"
	log "github.com/sirupsen/logrus"
	"sync"
	"test/config"
	"time"
)

var (
	pool *redis.Pool
)

//初始化一个pool
func newPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		MaxActive:   5,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", password); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

type Redis struct {
	sync.RWMutex
	Pool *redis.Pool
}

var RedisSource Redis

func InitRedis() {
	pool = newPool(config.Sysconfig.Redis.Addr, config.Sysconfig.Redis.Password)
	RedisSource.Pool = pool
	log.Print("连接池初始化成功")
}
func (r *Redis) RedisSet(key, value interface{}) bool {
	r.Lock()
	conn := r.Pool.Get()
	defer conn.Close()
	r.Unlock()
	_, err := conn.Do("set", key, value)
	if err != nil {
		return false
	}
	return true
}

//设置值并设置过期时间
func (r *Redis) RedisSetEx(key, value, time interface{}) bool {
	r.Lock()
	conn := r.Pool.Get()
	defer conn.Close()
	r.Unlock()
	_, err := conn.Do("SETEX", key, time, value)
	if err != nil {
		return false
	}
	return true
}
func (r *Redis) RedisGet(key string) string {
	r.Lock()
	conn := r.Pool.Get()
	defer conn.Close()
	r.Unlock()
	v, err := redis.String(conn.Do("get", key))
	if err != nil {
		return ""
	}
	return v
}

func (r *Redis) RedisSetBool(key string, value bool) error {
	r.Lock()
	conn := r.Pool.Get()
	defer conn.Close()
	r.Unlock()
	_, err := conn.Do("set", key, value)
	return err
}
func (r *Redis) RedisGetBool(key string) bool {
	r.Lock()
	conn := r.Pool.Get()
	defer conn.Close()
	r.Unlock()
	v, err := redis.Bool(conn.Do("get", key))
	if err != nil {
		return false
	}
	return v
}
func (r *Redis) RedisRemove(key interface{}) error {
	r.Lock()
	conn := r.Pool.Get()
	defer conn.Close()
	r.Unlock()
	_, err := conn.Do("del", key)
	return err
}
func (r *Redis) RedisExpire(key interface{}, time int) error {
	r.Lock()
	conn := r.Pool.Get()
	defer conn.Close()
	r.Unlock()
	_, err := conn.Do("EXPIRE", key, time)
	return err
}
