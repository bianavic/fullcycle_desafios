package internal

import "net/http"

type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

type RealHTTPClient struct {
	Client *http.Client
}

func (c *RealHTTPClient) Get(url string) (*http.Response, error) {
	return c.Client.Get(url)
}
