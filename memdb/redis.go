package memdb

import (
	"ddaom/define"
	"ddaom/domain"
	"fmt"
	"log"
	"os"

	"github.com/gomodule/redigo/redis"
)

var pool *redis.Pool

const (
	USER_TOKEN = "USER:TOKEN:"
)

func initRedis() {
	// init redis connection pool
	initPool()

	// bootstramp some data to redis
	// initStore()
}

func initPool() {
	pool = &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,

		Dial: func() (redis.Conn, error) {
			conn, err := redis.DialURL(define.DSN_REDIS, redis.DialDatabase(1))
			if err != nil {
				log.Printf("ERROR: fail init redis: %s", err.Error())
				os.Exit(1)
			}
			return conn, err
		},
	}
}

func Ping(conn redis.Conn) {
	_, err := redis.String(conn.Do("PING"))
	if err != nil {
		log.Printf("ERROR: fail ping redis conn: %s", err.Error())
		os.Exit(1)
	}
}

func Zadd(key string, score int, data interface{}) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("ZADD", key, score, data)
	if err != nil {
		log.Printf("ERROR: fail set key %s, score %d, error %s", key, score, err.Error())
		return err
	}

	return nil
}

func Set(key string, val string) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, val)
	if err != nil {
		log.Printf("ERROR: fail set key %s, val %s, error %s", key, val, err.Error())
		return err
	}

	return nil
}

func Get(key string) (string, error) {
	conn := pool.Get()
	defer conn.Close()

	s, err := redis.String(conn.Do("GET", key))
	if err != nil {
		log.Printf("ERROR: fail get key %s, error %s", key, err.Error())
		return "", err
	}

	return s, nil
}

func IsExistSession(token string) bool {
	conn := pool.Get()
	defer conn.Close()

	ret, _ := redis.Values(conn.Do("HGETALL", USER_TOKEN+token))
	// fmt.Println(len(ret))
	if len(ret) > 0 {
		return true
	} else {
		return false
	}
}

func SetSession(token string, userToken domain.UserToken) error {
	conn := pool.Get()
	defer conn.Close()

	// set date
	_, err := conn.Do("HMSET", redis.Args{}.Add(USER_TOKEN+token).AddFlat(userToken)...)
	return err
}

func SetSessionExpireAdd(token string) error {
	conn := pool.Get()
	defer conn.Close()
	_, err := conn.Do("EXPIRE", USER_TOKEN+token, 60*15)
	return err
}

func Sadd(key string, val string) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("SADD", key, val)
	if err != nil {
		log.Printf("ERROR: fail add val %s to set %s, error %s", val, key, err.Error())
		return err
	}

	return nil
}

func Smembers(key string) ([]string, error) {
	conn := pool.Get()
	defer conn.Close()

	s, err := redis.Strings(conn.Do("SMEMBERS", key))
	if err != nil {
		log.Printf("ERROR: fail get set %s , error %s", key, err.Error())
		return nil, err
	}

	return s, nil
}

func HVALS(key string) ([]string, error) {
	conn := pool.Get()
	defer conn.Close()
	s, err := redis.Strings(conn.Do("HVALS", key))
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return s, nil
}

func HMGET(values ...interface{}) ([]string, error) {
	conn := pool.Get()
	defer conn.Close()
	// fmt.Println(values...)
	s, err := redis.Strings(conn.Do("HMGET", values...))
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return s, nil
}

func HGETALL(value string) ([]string, error) {
	conn := pool.Get()
	defer conn.Close()
	// fmt.Println(value)
	s, err := redis.Strings(conn.Do("HGETALL", value))
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return s, nil
}

func ZMSCORE(values ...interface{}) ([]string, error) {
	conn := pool.Get()
	defer conn.Close()
	s, err := redis.Strings(conn.Do("ZMSCORE", values...))
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return s, nil
}

func ZREVRANGE(key string) ([]string, error) {
	conn := pool.Get()
	defer conn.Close()
	s, err := redis.Strings(conn.Do("ZREVRANGE", key, 0, -1, "WITHSCORES"))
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return s, nil
}

func Execute(exec string, values ...interface{}) ([]string, error) {
	conn := pool.Get()
	defer conn.Close()
	s, err := redis.Strings(conn.Do(exec, values...))
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return s, nil
}

func RunRedis() {
	initRedis()
}
