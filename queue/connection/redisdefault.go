package connection

import (
	"github.com/go-redis/redis"
	"time"
	"windcontrol-go/config/types"
	redisSaver "windcontrol-go/persist/redis"
)

var redisClient *redis.Client

var redisDefaultBlockTime = time.Second * 2

var redisDefaultQueueName = "redis_default"

func RedisDefault() (types.Task, error) {
	var err error
	if redisClient == nil {
		redisClient, err = redisSaver.NewClient()
		if err != nil {
			panic(err)
		}
	}

	list, err := redisClient.BLPop(redisDefaultBlockTime, redisDefaultQueueName).Result()
	if err != nil {
		return nil, err
	}

	result := ""
	if len(list) > 0 {
		result = list[0]
	}

	return result, err
}