package pokeapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	http    *http.Client
	baseURL string
}

func NewClient(baseURL string) *Client {
	return &Client{
		http: &http.Client{
			Timeout: 5 * time.Second,
		},
		baseURL: baseURL,
	}
}

func (c *Client) get(ctx context.Context, path string, response any) error {
	url := fmt.Sprintf("%s%s", c.baseURL, path)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}

	return json.NewDecoder(resp.Body).Decode(response)
}
