package pkg

func ConvertTemperature(celsius float64) map[string]float64 {
	return map[string]float64{
		"temp_C": celsius,
		"temp_F": celsius*1.8 + 32,
		"temp_K": celsius + 273.15,
	}
}
