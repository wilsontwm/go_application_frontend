package utils

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	Pool *redis.Pool
)

func init() {
	redisPort := os.Getenv("redis_port")
	if redisPort == "" {
		redisPort = "6379"
	}

	Pool = redisPool(redisPort)
	cleanupHook()
}

// Initialize the redis pool
func redisPool(port string) *redis.Pool {
	redisPool := redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":"+port)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return &redisPool
}

func cleanupHook() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)
	go func() {
		<-c
		Pool.Close()
		os.Exit(0)
	}()
}

// Ping the redis server
func RedisPing() error {

	conn := Pool.Get()
	defer conn.Close()

	_, err := redis.String(conn.Do("PING"))
	if err != nil {
		return fmt.Errorf("Cannot 'PING' db: %v", err)
	}
	return nil
}

// Get the value by passing in the key in redis
func RedisGet(key string) ([]byte, error) {

	conn := Pool.Get()
	defer conn.Close()

	var data []byte
	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return data, fmt.Errorf("Error getting key %s: %v", key, err)
	}
	return data, err
}

// Set the value in redis
func RedisSet(key string, value []byte) error {

	conn := Pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	if err != nil {
		v := string(value)
		if len(v) > 15 {
			v = v[0:12] + "..."
		}
		return fmt.Errorf("Error setting key %s to %s: %v", key, v, err)
	}
	return err
}

// Check if key esists in redis
func RedisExists(key string) (bool, error) {

	conn := Pool.Get()
	defer conn.Close()

	ok, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return ok, fmt.Errorf("Error checking if key %s exists: %v", key, err)
	}
	return ok, err
}

// Remove the key in redis
func RedisDelete(key string) error {

	conn := Pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	return err
}

// Get all the keys with pattern
func RedisGetKeys(pattern string) ([]string, error) {

	conn := Pool.Get()
	defer conn.Close()

	iter := 0
	keys := []string{}
	for {
		arr, err := redis.Values(conn.Do("SCAN", iter, "MATCH", pattern))
		if err != nil {
			return keys, fmt.Errorf("Error retrieving '%s' keys", pattern)
		}

		iter, _ = redis.Int(arr[0], nil)
		k, _ := redis.Strings(arr[1], nil)
		keys = append(keys, k...)

		if iter == 0 {
			break
		}
	}

	return keys, nil
}

// Increase counter in the redis
func RedisIncr(counterKey string) (int, error) {

	conn := Pool.Get()
	defer conn.Close()

	return redis.Int(conn.Do("INCR", counterKey))
}
