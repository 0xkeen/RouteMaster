package db

import "github.com/opencensus-integrations/redigo/redis"

type RedisDB struct {
	redis *redis.Pool
}

func NewRedisDB(pool *redis.Pool) *RedisDB {
	return &RedisDB{redis: pool}
}
