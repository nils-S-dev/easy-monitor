package api

import (
	"easy-monitor/internal/config"
	"easy-monitor/internal/monitor"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	// API versioning
	r.Route("/api/v1", func(r chi.Router) {
		r.Mount("/", Routes())
		r.Mount("/config", config.Routes())
		r.Mount("/monitors", monitor.Routes())
	})

	return r
}

func Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/health", healthHandler)

	return r
}
