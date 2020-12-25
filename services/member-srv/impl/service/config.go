package service

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

func configParms() (parms string, err error) {
	return
}

func casheCommet(conf string) (pool *redis.Pool, err error) {
	p := &redis.Pool{
		MaxIdle:     32,
		MaxActive:   100,
		IdleTimeout: time.Duration(100) * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				"127.0.0.1:6379",
				redis.DialPassword("dader487"),
			)
		},
	}
	return p, nil
}
