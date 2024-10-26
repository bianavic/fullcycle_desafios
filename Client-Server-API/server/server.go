package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Cotacao struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

type Cotacoes struct {
	USDBRL Cotacao `json:"usdbrl"`
}

func main() {

	http.HandleFunc("/", BuscaCotacaoHandler)
	http.ListenAndServe(":8080", nil)
}

func BuscaCotacaoHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	moedaParam := r.URL.Query().Get("moedas")
	if moedaParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	moeda, err := BuscaCotacao(moedaParam)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(moeda)
}

func BuscaCotacao(moedas string) (*Cotacoes, error) {
	resp, err := http.Get("https://economia.awesomeapi.com.br/json/last/" + moedas)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var c Cotacoes
	err = json.Unmarshal(body, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
