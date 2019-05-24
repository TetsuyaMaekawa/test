package redis

import (
	"log"

	"github.com/garyburd/redigo/redis"
)

// RedisConnection ...
func RedisConnection() redis.Conn {
	const IPPotr = "localhost:6379"

	redisConn, err := redis.Dial("tcp", IPPotr)
	if err != nil {
		log.Print(err)
	}
	return redisConn
}
