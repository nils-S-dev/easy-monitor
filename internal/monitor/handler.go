package monitor

import (
	"easy-monitor/internal/config"
	"encoding/json"
	"net/http"
)

func getMonitorsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	conf := config.GetConfig()
	data := GetMonitorResults(conf.Monitors)
	json.NewEncoder(w).Encode(data)
}
