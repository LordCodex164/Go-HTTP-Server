package handlers

import (
	"encoding/json"
	"net/http"
	"github.com/LordCodex164/httpserver/internal/metrics"
)

// Metrics handler exposes application metrics
func Metrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	snapshot := metrics.GetInstance().GetSnapshot()

	response := map[string]interface{}{
		"requests": map[string]interface{}{
			"total":   snapshot.TotalRequests,
			"success": snapshot.SuccessRequests,
			"client_errors": snapshot.ClientErrors,
			"server_errors": snapshot.ServerErrors,
		},
		"latency": map[string]interface{}{
			"avg_ms": snapshot.AvgLatency.Milliseconds(),
			"min_ms": snapshot.MinLatency.Milliseconds(),
			"max_ms": snapshot.MaxLatency.Milliseconds(),
		},
		"status_codes": snapshot.StatusCodes,
		"endpoints":    snapshot.EndpointCounts,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}