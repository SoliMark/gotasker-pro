package db

import (
	"context"
	"time"

	redis "github.com/redis/go-redis/v9"
)

func NewRedisClient(addr, pwd string, database int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,
		DB:       database,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = rdb.Ping(ctx).Err()

	return rdb
}
