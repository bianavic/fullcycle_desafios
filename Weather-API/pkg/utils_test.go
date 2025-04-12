package pkg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvertTemperature(t *testing.T) {

	t.Run("successfully convert temperature for positive celsius", func(t *testing.T) {
		result := ConvertTemperature(25.0)

		expected := map[string]float64{
			"temp_C": 25.0,
			"temp_F": 77.0,
			"temp_K": 298.15,
		}

		assert.Equal(t, expected["temp_C"], result["temp_C"], "temp_C mismatch")
		assert.Equal(t, expected["temp_F"], result["temp_F"], "temp_F mismatch")
		assert.Equal(t, expected["temp_K"], result["temp_K"], "temp_K mismatch")
	})

	t.Run("successfully convert temperature for zero celsius", func(t *testing.T) {
		result := ConvertTemperature(0.0)

		expected := map[string]float64{
			"temp_C": 0.0,
			"temp_F": 32.0,
			"temp_K": 273.15,
		}

		assert.Equal(t, expected["temp_C"], result["temp_C"], "temp_C mismatch")
		assert.Equal(t, expected["temp_F"], result["temp_F"], "temp_F mismatch")
		assert.Equal(t, expected["temp_K"], result["temp_K"], "temp_K mismatch")
	})

	t.Run("successfully convert temperature for negative celsius", func(t *testing.T) {
		result := ConvertTemperature(-10.0)

		expected := map[string]float64{
			"temp_C": -10.0,
			"temp_F": 14.0,
			"temp_K": 263.15,
		}

		assert.Equal(t, expected["temp_C"], result["temp_C"], "temp_C mismatch")
		assert.Equal(t, expected["temp_F"], result["temp_F"], "temp_F mismatch")
		assert.Equal(t, expected["temp_K"], result["temp_K"], "temp_K mismatch")
	})
}
