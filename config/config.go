package config

import (
	"log"
	"time"

	"github.com/joeshaw/envdecode"
)

type AppConf struct {
	Server     ServerConf
	RedisHosts RedisConf
}

type ServerConf struct {
	Port           int           `env:"SERVER_PORT,required"`
	TimeoutRead    time.Duration `env:"SERVER_TIMEOUT_READ,required"`
	TimeoutWrite   time.Duration `env:"SERVER_TIMEOUT_WRITE,required"`
	TimeoutIdle    time.Duration `env:"SERVER_TIMEOUT_IDLE,required"`
	HandlerTimeout time.Duration `env:"SERVER_HANDLER_TIMEOUT,required"`
}

type RedisConf struct {
	Host     string `env:"REDIS_HOST,required"`
	Port     int    `env:"REDIS_PORT,required"`
	Username string `env:"REDIS_USERNAME,required"`
	Password string `env:"REDIS_PASSWORD,required"`
}

func AppConfig() *AppConf {
	var c AppConf
	if err := envdecode.StrictDecode(&c); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}

	return &c
}
