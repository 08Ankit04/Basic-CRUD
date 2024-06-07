package router

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"github.com/Basic-CRUD/app/server"
)

func New(srv *server.Server) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

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

		r.Get("/employee", srv.HandleListEmployee)
		r.Get("/employee/{id}", srv.HandleGetEmployee)
		r.Post("/employee", srv.HandleCreateEmployee)
		r.Put("/employee/{id}", srv.HandleUpdateEmployee)
		r.Delete("/employee/{id}", srv.HandleDeleteEmployee)
	})

	return r
}
