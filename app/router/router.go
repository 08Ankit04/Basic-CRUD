package router

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"github.com/Basic-CRUD/app/server"
)

func New(srv *server.Server, hd time.Duration, hdw time.Duration) *chi.Mux {
	r := chi.NewRouter()
	l := srv.Logger

	l.Print("")

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowedHeaders:   []string{"Accept", "Content-Type", "X-CSRF-Token", "pragma", "Referer"},
		AllowCredentials: true,
		ExposedHeaders:   []string{"*"},
		MaxAge:           300,
	})

	r.Use(corsHandler.Handler)

	r.Get("/healthz", server.HandleLive)

	r.Route("/api/v1", func(r chi.Router) {

	})

	return r
}
