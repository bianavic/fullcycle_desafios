package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/bianavic/fullcycle_desafios/Weather-API-Otel/config"
	"github.com/bianavic/fullcycle_desafios/Weather-API-Otel/domain"
	"github.com/bianavic/fullcycle_desafios/Weather-API-Otel/dto"
)

func CEPHandler(w http.ResponseWriter, r *http.Request) {
	cfg := config.Load()

	var req dto.CEPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		domain.RespondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	// Forward to Service B
	resp, err := forwardToServiceB(cfg.ServiceBURL+"/weather", req.CEP)
	if err != nil {
		domain.RespondWithError(w, http.StatusInternalServerError, domain.ErrServiceBUnavailable.Error())
		return
	}
	defer resp.Body.Close()

	// Forward the response from Service B
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)

	if _, err := bytes.NewBuffer(nil).ReadFrom(resp.Body); err != nil {
		domain.RespondWithError(w, http.StatusInternalServerError, "failed to read response")
		return
	}
}

func forwardToServiceB(url, cep string) (*http.Response, error) {
	reqBody, _ := json.Marshal(map[string]string{"cep": cep})
	return http.Post(url, "application/json", bytes.NewBuffer(reqBody))
}
