package memdb

import (
	"ddaom/define"
	"ddaom/domain"
	"fmt"
	"log"
	"os"

	"github.com/gomodule/redigo/redis"
)

var poolMaster *redis.Pool
var poolSlave *redis.Pool

const (
	USER_TOKEN = "USER:TOKEN:"
)

func RunRedis() {
	initRedis()
}

func initRedis() {
	initPool()
}

func initPool() {
	poolMaster = &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.DialURL(define.Mconn.DsnRedisMaster, redis.DialDatabase(1))
			if err != nil {
				log.Printf("ERROR: fail init redis: %s", err.Error())
				os.Exit(1)
			}
			return conn, err
		},
	}
	poolSlave = &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.DialURL(define.Mconn.DsnRedisSlave, redis.DialDatabase(1))
			if err != nil {
				log.Printf("ERROR: fail init redis: %s", err.Error())
				os.Exit(1)
			}
			return conn, err
		},
	}
}

func Zadd(key string, score int, data interface{}) error {
	conn := poolMaster.Get()
	defer conn.Close()

	_, err := conn.Do("ZADD", key, score, data)
	if err != nil {
		log.Printf("ERROR: fail set key %s, score %d, error %s", key, score, err.Error())
		return err
	}

	return nil
}

func Set(key string, val string) error {
	conn := poolMaster.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, val)
	if err != nil {
		log.Printf("ERROR: fail set key %s, val %s, error %s", key, val, err.Error())
		return err
	}

	return nil
}

func Get(key string) (string, error) {
	conn := poolSlave.Get()
	defer conn.Close()

	s, err := redis.String(conn.Do("GET", key))
	if err != nil {
		log.Printf("ERROR: fail get key %s, error %s", key, err.Error())
		return "", err
	}

	return s, nil
}

func IsExistSession(token string) bool {
	conn := poolSlave.Get()
	defer conn.Close()

	ret, _ := redis.Values(conn.Do("HGETALL", USER_TOKEN+token))
	if len(ret) > 0 {
		return true
	} else {
		return false
	}
}

func SetSession(token string, userToken domain.UserToken) error {
	conn := poolMaster.Get()
	defer conn.Close()
	_, err := conn.Do("HMSET", redis.Args{}.Add(USER_TOKEN+token).AddFlat(userToken)...)
	return err
}

func SetSessionExpireAdd(token string) error {
	conn := poolSlave.Get()
	defer conn.Close()
	_, err := conn.Do("EXPIRE", USER_TOKEN+token, 60*15)
	return err
}

func Sadd(key string, val string) error {
	conn := poolMaster.Get()
	defer conn.Close()

	_, err := conn.Do("SADD", key, val)
	if err != nil {
		log.Printf("ERROR: fail add val %s to set %s, error %s", val, key, err.Error())
		return err
	}

	return nil
}

func ZMSCORE(values ...interface{}) ([]string, error) {
	conn := poolSlave.Get()
	defer conn.Close()
	s, err := redis.Strings(conn.Do("ZMSCORE", values...))
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return s, nil
}

func ZREVRANGE(key string) ([]string, error) {
	conn := poolSlave.Get()
	defer conn.Close()
	s, err := redis.Strings(conn.Do("ZREVRANGE", key, 0, -1, "WITHSCORES"))
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return s, nil
}
