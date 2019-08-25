package redis

import (
	"github.com/go-redis/redis"
	"windcontrol-go/config/persist"
)

func NewClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: persist.RedisAddr,
		Password: persist.RedisPassword,
		DB: persist.Db,
	})

	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}

	return client, nil
}