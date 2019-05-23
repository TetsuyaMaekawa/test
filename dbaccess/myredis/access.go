package myredis

import (
	"log"
	"time"

	"github.com/garyburd/redigo/redis"
)

// OpenRedis Redisとの接続を確立
func OpenRedis() (*redis.Pool, error) {

	// redis.PoolのDialがfunc()でエラーが取り出せないため処理を分割
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return conn, err
		},
	}, nil
}
