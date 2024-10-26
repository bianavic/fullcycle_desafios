package main

import (
	"io"
	"net/http"
)

func main() {
	req, err := http.Get("https://economia.awesomeapi.com.br/json")
	if err != nil {
		panic(err)
	}
	res, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	println(string(res))
}
