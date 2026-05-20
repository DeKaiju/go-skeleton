package redis

import (
	"fmt"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
)

var pool *redis.Pool

func GetCache() {
	host := os.Getenv("CACHE_HOST")
	port := os.Getenv("CACHE_PORT")
	pass := os.Getenv("CACHE_PASS")
	if host == "" {
		host = "127.0.0.1"
	}
	if port == "" {
		port = "6379"
	}
	server := fmt.Sprintf("%s:%s", host, port)
	pool = &redis.Pool{
		MaxActive:   100,
		MaxIdle:     100,
		IdleTimeout: 240 * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			// Authentication
			if pass != "" {
				if _, err := c.Do("AUTH", pass); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, nil
		},
	}
}

func Get() redis.Conn {
	if pool == nil {
		GetCache()
	}

	return pool.Get()
}

func Close() error {
	if pool == nil {
		return nil
	}

	err := pool.Close()
	pool = nil
	return err
}
