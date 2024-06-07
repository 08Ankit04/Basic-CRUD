package server

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
)

type Server struct {
	Redis     *redis.Client
	Logger    *log.Logger
	Validator *validator.Validate
}

func New(
	redis *redis.Client,
	logger *log.Logger,
	validator *validator.Validate,
) *Server {
	return &Server{
		Redis:     redis,
		Logger:    logger,
		Validator: validator,
	}
}
