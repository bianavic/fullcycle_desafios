package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMaskAPIKey(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("abcd****1234", maskAPIKey("abcd1234abcd1234"))
	assert.Equal("******", maskAPIKey("123"))
}

func TestGetServerPort(t *testing.T) {
	assert := assert.New(t)

	_ = os.Setenv("PORT", "9090")
	assert.Equal("9090", getServerPort())

	_ = os.Unsetenv("PORT")
	assert.Equal("8080", getServerPort())
}

func TestGetAPIKey(t *testing.T) {
	assert := assert.New(t)

	_ = os.Setenv("WEATHER_API_KEY", "testapikey123456")
	key := getAPIKey()
	assert.Equal("testapikey123456", key)

	_ = os.Unsetenv("WEATHER_API_KEY")
}

func TestLoadEnv(t *testing.T) {
	assert := assert.New(t)

	_ = os.Setenv("ENV", "production")
	err := loadEnv()
	assert.NoError(err)

	_ = os.Setenv("ENV", "development")
	err = loadEnv()
	assert.Error(err)
}
