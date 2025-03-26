package crawler

import (
	"context"
	"net/http"
	"time"
)

type HttpClient struct {
	Client        *http.Client
	Timeout       time.Duration
	ContextCancel context.CancelFunc
}

func NewHttpClient() *HttpClient {
	timeout := 30 * time.Second
	return &HttpClient{
		Client: &http.Client{
			Timeout: timeout,
		},
		Timeout: timeout,
	}
}

func (hc *HttpClient) Do(request *http.Request) (*http.Response, error) {
	return hc.Client.Do(request)
}

func (hc *HttpClient) NewGet(url string) (*http.Request, error) {
	ctx, cancel := context.WithTimeout(context.Background(), hc.Timeout)
	hc.ContextCancel = cancel
	return http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
}

func (hc *HttpClient) NewPost(url string) (*http.Request, error) {
	ctx, cancel := context.WithTimeout(context.Background(), hc.Timeout)
	hc.ContextCancel = cancel
	return http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
}

func (hc *HttpClient) Close() {
	hc.ContextCancel()
	hc.Client.CloseIdleConnections()

}
