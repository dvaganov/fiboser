package main

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	redisDBParam       = "db"
	redisPasswordParam = "password"
)

func NewRedis(dsn string) (*redis.Client, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	db, err := strconv.Atoi(u.Query().Get(redisDBParam))
	if err != nil {
		return nil, err
	}

	fmt.Println(u.Host)
	rdb := redis.NewClient(&redis.Options{
		Addr:     u.Host,
		Password: u.Query().Get(redisPasswordParam),
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}
