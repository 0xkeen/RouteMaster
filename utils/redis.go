package utils

import (
	"time"

	"github.com/opencensus-integrations/redigo/redis"
)

func GetRedisPool(host, password string) (*redis.Pool, error) {
	RedisPool := &redis.Pool{
		MaxIdle:     10,
		MaxActive:   10000,
		IdleTimeout: 60 * time.Second,
		Dial: func() (redis.Conn, error) {
			var err error
			var c redis.Conn
			if password == "" {
				c, err = redis.DialURL(host)
			} else {
				c, err = redis.DialURL(host, redis.DialPassword(password))
			}
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
	return RedisPool, nil
}
