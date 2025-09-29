package config

import (
	"encoding/json"
	"net/http"
)

func getConfigHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	config := GetConfig()
	json.NewEncoder(w).Encode(config)
}
