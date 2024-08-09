package extension

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
)

type Redis struct {
	redis *redis.Pool
}

func (m *Redis) Init() error {
	if m.redis != nil {
		return nil
	}
	database := os.Getenv("REDIS_DATABASE")
	password := os.Getenv("REDIS_PASSWORD")
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	maxIdleConns := 2
	maxOpenConns := 3000
	maxLifeTime := 24
	readTimeout := 10
	writeTimeout := 10
	if len(os.Getenv("REDIS_MAX_IDLE_CONNS")) > 0 {
		maxIdleConns, _ = strconv.Atoi(os.Getenv("REDIS_MAX_IDLE_CONN"))
	}
	if len(os.Getenv("REDIS_MAX_OPEN_CONNS")) > 0 {
		maxOpenConns, _ = strconv.Atoi(os.Getenv("REDIS_MAX_OPEN_CONNS"))
	}
	if len(os.Getenv("REDIS_MAX_LIFE_TIME")) > 0 {
		maxLifeTime, _ = strconv.Atoi(os.Getenv("REDIS_MAX_LIFE_TIME"))
	}
	if len(os.Getenv("REDIS_READ_TIMEOUT")) > 0 {
		readTimeout, _ = strconv.Atoi(os.Getenv("REDIS_READ_TIMEOUT"))
	}
	if len(os.Getenv("REDIS_WRITE_TIMEOUT")) > 0 {
		writeTimeout, _ = strconv.Atoi(os.Getenv("REDIS_WRITE_TIMEOUT"))
	}

	address := fmt.Sprintf("%s:%s", host, port)
	db, _ := strconv.Atoi(database)

	m.redis = &redis.Pool{
		MaxIdle:   maxIdleConns, // idle connect num
		MaxActive: maxOpenConns, // max connect num
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				address,
				redis.DialPassword(password),
				redis.DialDatabase(db),
				redis.DialWriteTimeout(time.Duration(readTimeout)*time.Second),
				redis.DialReadTimeout(time.Duration(writeTimeout)*time.Second),
			)
		},
		IdleTimeout: time.Duration(maxLifeTime) * time.Second,
		Wait:        true,
	}

	return m.testConnect()
}

func (m *Redis) Redis() *redis.Pool {
	return m.redis
}

func (m *Redis) RedisConnect() redis.Conn {
	conn := m.redis.Get()
	return conn
}

func (m *Redis) RedisExec(cmd string, params ...any) (reply interface{}, err error) {
	conn := m.RedisConnect()
	reply, err = conn.Do(cmd, params...)
	conn.Close()
	return reply, err
}

func (m *Redis) RedisGet(key string) (reply interface{}, err error) {
	return m.RedisExec("get", key)
}

func (m *Redis) RedisSet(key string, value any, ttl ...int) (reply interface{}, err error) {
	if len(ttl) > 0 {
		return m.RedisExec("set", key, value, "EX", ttl[0])
	} else {
		return m.RedisExec("set", key, value)
	}
}

func (m *Redis) RedisDelete(key string) (reply interface{}, err error) {
	return m.RedisExec("del", key)
}

func (m *Redis) testConnect() error {
	conn := m.RedisConnect()
	res, err := conn.Do("ping")
	log.Println("redis test: ping =", res)
	return err
}
