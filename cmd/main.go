package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Basic-CRUD/adapter/redis"
	"github.com/Basic-CRUD/app/router"
	"github.com/Basic-CRUD/app/server"
	"github.com/go-playground/validator/v10"

	"github.com/Basic-CRUD/config"
)

// @title			Basic Crud service API
// @version		0.1
// @description	This is the API documentation of the basic-crud service
// @license.name	MIT
// @host			localhost:443
// @basePath		/api/v1
// @schemes		http https
func main() {
	// Get configurations
	appConf := config.AppConfig()

	// Get logger
	appLogger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)

	// Get Redis connection
	appCache := redis.New(&appConf.RedisHosts)
	redisWrapper := server.NewRedisWrapper(appCache)

	

	// Get validator
	appValidator := validator.New()

	// Set server dependencies
	srv := server.New(redisWrapper, appLogger, appValidator)

	// Get server routes
	appRouter := router.New(srv)

	// Create server
	address := fmt.Sprintf(":%d", appConf.Server.Port)
	appLogger.Printf("Starting server %v", address)
	s := &http.Server{
		Addr:         address,
		Handler:      appRouter,
		ReadTimeout:  appConf.Server.TimeoutRead,
		WriteTimeout: appConf.Server.TimeoutWrite,
		IdleTimeout:  appConf.Server.TimeoutIdle,
	}

	closed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)
		<-sigint

		appLogger.Print("Shutting down server...")

		if err := s.Shutdown(context.Background()); err != nil {
			appLogger.Fatalf("Server shutdown failed: %v", err)
		}

		close(closed)
	}()

	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		appLogger.Fatalf("Server startup failed: %v", err)
	}

	<-closed
}
