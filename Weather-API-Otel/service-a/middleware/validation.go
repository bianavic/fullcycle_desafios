package middleware

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"net/http"
	"regexp"

	"github.com/bianavic/fullcycle_desafios/Weather-API-Otel/domain"
	"github.com/bianavic/fullcycle_desafios/Weather-API-Otel/dto"
)

var validate = validator.New()

func ValidateCEP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req dto.CEPRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			RespondWithError(w, http.StatusBadRequest, "invalid request")
			return
		}

		if err := validate.Struct(req); err != nil {
			RespondWithError(w, http.StatusUnprocessableEntity, domain.ErrInvalidCEP.Error())
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isValidCEP(cep string) bool {
	return regexp.MustCompile(`^\d{8}$`).MatchString(cep)
}
