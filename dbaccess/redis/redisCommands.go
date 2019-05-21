package redis

import (
	"log"

	"github.com/garyburd/redigo/redis"
)

// RedisSet string型のデータを期限付きで設定
func RedisSet(key string, value string, second int, c redis.Conn) {
	c.Do("SETEX", key, second, value)
}

// RedisGet ...
func RedisGet(key string, c redis.Conn) string {
	str, err := redis.String(c.Do("GET", key))
	if err != nil {
		log.Print(err)
	}
	return str
}
