package handlers

import "github.com/stretchr/testify/mock"

type MockWeatherUsecase struct {
	mock.Mock
}

func (m *MockWeatherUsecase) GetWeatherByCEP(cep string) (interface{}, error) {
	args := m.Called(cep)
	return args.Get(0), args.Error(1)
}
