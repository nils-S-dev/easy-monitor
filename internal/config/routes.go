package config

import "github.com/go-chi/chi/v5"

func Routes() chi.Router {
	var r = chi.NewRouter()
	r.Get("/", getConfigHandler)
	return r
}
