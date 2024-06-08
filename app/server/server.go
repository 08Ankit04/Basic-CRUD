package server

import (
	"context"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
)

type Server struct {
	Redis     RedisClient
	Logger    *log.Logger
	Validator *validator.Validate
}

func New(
	redis RedisClient,
	logger *log.Logger,
	validator *validator.Validate,
) *Server {
	return &Server{
		Redis:     redis,
		Logger:    logger,
		Validator: validator,
	}
}

type RedisClient interface {
	Get(key string) *redis.StringCmd
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Del(key string) *redis.IntCmd
	Keys(pattern string) *redis.StringSliceCmd
}

type RedisWrapper struct {
	client *redis.Client
}

func NewRedisWrapper(client *redis.Client) *RedisWrapper {
	return &RedisWrapper{client: client}
}

func (rw *RedisWrapper) Get(key string) *redis.StringCmd {
	return rw.client.Get(context.Background(), key)
}

func (rw *RedisWrapper) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return rw.client.Set(context.Background(), key, value, expiration)
}

func (rw *RedisWrapper) Del(key string) *redis.IntCmd {
	return rw.client.Del(context.Background(), key)
}

func (rw *RedisWrapper) Keys(pattern string) *redis.StringSliceCmd {
	return rw.client.Keys(context.Background(), pattern)
}
