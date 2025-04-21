package domain

import (
	"encoding/json"
	"errors"
	"github.com/bianavic/fullcycle_desafios/Weather-API-Otel/dto"
	"net/http"
)

var (
	ErrInvalidCEP          = errors.New("invalid zipcode")
	ErrCEPNotFound         = errors.New("can not find zipcode")
	ErrServiceBUnavailable = errors.New("service b unavailable")
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

type ValidationError struct {
	message string
}

func (e *ValidationError) Error() string {
	return e.message
}

func NewValidationError(message string) error {
	return &ValidationError{message: message}
}

// WriteError writes a standardized error response
func WriteError(w http.ResponseWriter, statusCode int, message string, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ErrorResponse{
		Message: message,
	}

	if err != nil {
		response.Error = err.Error()
	}

	json.NewEncoder(w).Encode(response)
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(dto.ErrorResponse{Error: message})
}
