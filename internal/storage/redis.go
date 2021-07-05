package storage

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

func NewRedis(addr string) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:        addr,
		DB:          0,
		DialTimeout: time.Second * 10,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}
