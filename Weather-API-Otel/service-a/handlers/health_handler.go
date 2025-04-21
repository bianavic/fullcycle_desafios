package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bianavic/fullcycle_desafios/Weather-API-Otel/dto"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.HealthResponse{Status: "ok"})
}
