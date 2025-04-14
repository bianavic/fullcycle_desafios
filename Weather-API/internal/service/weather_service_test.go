package service

import (
	"encoding/json"
	"github.com/bianavic/fullcycle_desafios/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWeatherAPIService_GetWeatherByCity(t *testing.T) {
	const testAPIKey = "test-api-key"

	t.Run("successfully retrieves weather data for valid city", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v1/current.json", r.URL.Path)
			assert.Equal(t, "test-api-key", r.URL.Query().Get("key"))
			assert.Equal(t, "São Paulo", r.URL.Query().Get("q"))

			response := domain.WeatherAPIResponse{
				Location: struct {
					Name    string `json:"name"`
					Region  string `json:"region"`
					Country string `json:"country"`
				}{
					Name:    "São Paulo",
					Region:  "SP",
					Country: "Brazil",
				},
				Current: struct {
					TempC     float64 `json:"temp_c"`
					TempF     float64 `json:"temp_f"`
					Condition struct {
						Text string `json:"text"`
					} `json:"condition"`
				}{
					TempC: 25.5,
					TempF: 77.9,
					Condition: struct {
						Text string `json:"text"`
					}{
						Text: "Sunny",
					},
				},
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		}))
		defer ts.Close()

		service := NewWeatherAPIService(testAPIKey)
		service.BaseURL = ts.URL

		resp, err := service.GetWeatherByCity("São Paulo")
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 25.5, resp.Current.TempC)
	})

	t.Run("should return error when city is empty", func(t *testing.T) {
		service := NewWeatherAPIService("valid-api-key")

		resp, err := service.GetWeatherByCity("")
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "city cannot be empty")
	})

	t.Run("should return error when fail to send request", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			hj, ok := w.(http.Hijacker)
			if !ok {
				panic("server doesn't support hijacking")
			}
			conn, _, err := hj.Hijack()
			if err != nil {
				panic(err)
			}
			conn.Close()
		}))
		defer ts.Close()

		service := NewWeatherAPIService(testAPIKey)
		service.BaseURL = ts.URL

		_, err := service.GetWeatherByCity("London")

		assert.ErrorIs(t, err, domain.ErrFailedToSendRequest)
	})

	t.Run("should return error when fail to read JSON", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			// Hijack the connection to simulate a broken body
			hj, ok := w.(http.Hijacker)
			if !ok {
				panic("server doesn't support hijacking")
			}
			conn, buf, err := hj.Hijack()
			if err != nil {
				panic(err)
			}

			// Write valid headers but close connection before sending body
			buf.WriteString("HTTP/1.1 200 OK\r\n")
			buf.WriteString("Content-Type: application/json\r\n")
			buf.WriteString("\r\n")
			buf.Flush()
			conn.Close()
		}))
		defer ts.Close()

		service := NewWeatherAPIService(testAPIKey)
		service.BaseURL = ts.URL

		_, err := service.GetWeatherByCity("London")

		assert.ErrorIs(t, err, domain.ErrFailedToParseData)
	})

	t.Run("should return service unavailable when API returns HTML error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte("<html>Error</html>"))
		}))
		defer ts.Close()

		service := NewWeatherAPIService(testAPIKey)
		service.BaseURL = ts.URL

		resp, err := service.GetWeatherByCity("São Paulo")
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "HTML instead of JSON")
	})

	t.Run("should return error when fail to unmarshalling JSON", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("invalid json"))
		}))
		defer ts.Close()

		service := NewWeatherAPIService(testAPIKey)
		service.BaseURL = ts.URL

		resp, err := service.GetWeatherByCity("São Paulo")
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), domain.ErrFailedToParseData.Error())
	})
}
