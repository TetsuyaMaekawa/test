package myredis

import (
	"log"

	"github.com/garyburd/redigo/redis"
)

// SetKeyValue expireを30秒とし値をセット
func SetKeyValue(key, value string, RD *redis.Pool) {
	conn := RD.Get()
	conn.Do("SETEX", key, 30, value)
}

// GetKey キーを取得
func GetKey(key string, RD *redis.Pool) string {
	conn := RD.Get()
	rtnStr, err := redis.Strings(conn.Do("KEYS", key))
	if err != nil {
		log.Print(err)
	}
	return rtnStr[0]
}

// GetValue バリューを取得
func GetValue(key string, RD *redis.Pool) string {
	conn := RD.Get()
	rtnStr, err := redis.String(conn.Do("GET", key))
	if err != nil {
		log.Print(err)
	}
	return rtnStr
}
