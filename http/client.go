package http

import (
	"bytes"
	"net/http"
	"time"
)

var httpClient *HttpClient

type HttpClient struct {
	client *http.Client
}

func init() {
	// logger
}

func New() *HttpClient {
	if httpClient != nil {
		return httpClient
	}

	httpClient = &HttpClient{
		client: &http.Client{
			Timeout: 5 * time.Second,
			Transport: &http.Transport{
				IdleConnTimeout: 10 * time.Second,
				//This can be set up to avoid server cert verification or to add client certificates
				//TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}

	return httpClient
}

func (hc *HttpClient) MakeRequest(method, uri string, body []byte, headers map[string]string) (*http.Response, error) {
	var (
		request *http.Request
		err     error
	)

	if body != nil {
		request, err = http.NewRequest(
			method,
			uri,
			bytes.NewReader(body),
		)
	} else {
		request, err = http.NewRequest(
			method,
			uri,
			nil,
		)
	}

	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		request.Header.Set(k, v)
	}

	return hc.client.Do(request)
}
