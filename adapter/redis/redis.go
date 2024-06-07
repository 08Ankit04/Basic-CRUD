package redis

import (
	"fmt"

	"github.com/Basic-CRUD/config"

	"github.com/go-redis/redis/v8"
)

func New(conf *config.RedisConf) *redis.Client {
	options := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password: conf.Password,
		DB:       0,
	}

	return redis.NewClient(options)
}
