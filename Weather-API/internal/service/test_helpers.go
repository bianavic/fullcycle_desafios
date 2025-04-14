package service

import (
	"errors"
	"io"
	"net/http"
	"strings"
)

type MockHTTPClient struct {
	ResponseBody       string
	ResponseBodyReader io.ReadCloser
	StatusCode         int
	Err                error
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	body := m.ResponseBodyReader
	if body == nil {
		body = io.NopCloser(strings.NewReader(m.ResponseBody))
	}

	return &http.Response{
		StatusCode: m.StatusCode,
		Body:       body,
		Header:     make(http.Header),
	}, nil
}

type brokenReader struct{}

func (b brokenReader) Read(p []byte) (int, error) {
	return 0, errors.New("simulated read error")
}

func (b brokenReader) Close() error {
	return nil
}
