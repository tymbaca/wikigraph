package httpclient

import (
	"context"
	"net/http"

	"golang.org/x/time/rate"
)

type RateLimitingClient struct {
	client Client
	lim    *rate.Limiter
}

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewRateLimitingClient(client Client, rps, burst int) *RateLimitingClient {
	return &RateLimitingClient{
		client: client,
		lim:    rate.NewLimiter(rate.Limit(rps), burst),
	}
}

func (c *RateLimitingClient) Do(req *http.Request) (*http.Response, error) {
	if err := c.lim.Wait(context.TODO()); err != nil {
		return nil, err
	}

	return c.client.Do(req)
}
