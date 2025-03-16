package handlers

import (
	"context"
	"fmt"
	"github.com/bianavic/fullcycle_desafios/infra/api"
)

type Result struct {
	Source string
	Data   interface{}
	Error  error
}

type FetchAddressFunc func(ctx context.Context, cep, source string) (interface{}, error)

func FetchAddressHandler(ctx context.Context, cep, source string, fetchAddressFunc FetchAddressFunc, resultChan chan Result) {
	address, err := fetchAddressFunc(ctx, cep, source)
	if err != nil {
		resultChan <- Result{Source: source, Error: err}
		return
	}
	resultChan <- Result{Source: source, Data: address}
}

func PrintResult(result Result) {
	switch data := result.Data.(type) {
	case *api.ViaCepAPIResponse:
		fmt.Printf("Received from viacep: source:%s - %s, %s - %s, %s, %s\n", result.Source, data.Logradouro, data.Bairro, data.Localidade, data.Estado, data.Cep)
	case *api.BrasilAPIResponse:
		fmt.Printf("Received from brasilapi: source:%s - %s, %s - %s, %s, %s\n", result.Source, data.Street, data.Neighborhood, data.City, data.State, data.Cep)
	}
}
