package redis

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

// // RedisConnection ...
// func RedisConnection() redis.Conn {
// 	const IPPotr = "localhost:6379"

// 	redisConn, err := redis.Dial("tcp", IPPotr)
// 	if err != nil {
// 		log.Print(err)
// 	}
// 	return redisConn
// }

// NewPool redisのコネクションプールの生成
func NewPool() *redis.Pool {
	return &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
		MaxIdle:     3,
		MaxActive:   3,
		IdleTimeout: 240 * time.Second,
	}
}
