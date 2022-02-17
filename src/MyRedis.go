package src

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

var RedisDefaultPool *redis.Pool

func newPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", addr)
		},
	}
}
func init() {
	RedisDefaultPool = newPool("180.184.74.5:6379")
}
