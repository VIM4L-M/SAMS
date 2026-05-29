package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var startTime = time.Now()

func Health(rulesCount int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		uptime := time.Since(startTime)
		h := int(uptime.Hours())
		m := int(uptime.Minutes()) % 60
		s := int(uptime.Seconds()) % 60

		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":      "ok",
			"version":     "1.0.0",
			"rulesLoaded": rulesCount,
			"uptime":      fmt.Sprintf("%dh %dm %ds", h, m, s),
		})
	}
}
