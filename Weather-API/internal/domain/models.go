package domain

type LocationResponse struct {
	City     string `json:"city"`
	State    string `json:"state"`
	CEP      string `json:"cep"`
	District string `json:"district"`
	Street   string `json:"street"`
	Service  string `json:"service"`
}

type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

type TemperatureResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}
