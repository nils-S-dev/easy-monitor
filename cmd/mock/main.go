package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ApiResponse struct {
	Data interface{} `json:"data,omitempty"`
}

func main() {
	r := chi.NewRouter()

	beer1 := map[string]string{
		"name": "Augustiner",
	}

	beer2 := map[string]string{
		"name": "Berg",
	}

	beers := []map[string]string{
		beer1,
		beer2,
	}

	r.Get("/beer", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, ApiResponse{Data: beers})
	})

	r.Get("/beer/1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-hello-world", "Hello World!")
		writeJSON(w, http.StatusOK, ApiResponse{Data: beer1})
	})

	r.Get("/beer/-1", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusNotFound, ApiResponse{})
	})

	r.Post("/beer", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-hello", "World!")
		w.Header().Set("X-world", "Hello")
		writeJSON(w, http.StatusOK, ApiResponse{Data: beer2})
	})

	r.Delete("/beer/1", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusNoContent, ApiResponse{})
	})

	http.ListenAndServe(":8081", r)
}

// Helper: write JSON with status code
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
