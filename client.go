package pokeapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
	// a request to get can provide both an absolute and relative path. We conveniently
	// check if the path is absolute or relative and combine them if it's relative.
	parsed, err := url.Parse(path)
	if err != nil {
		return fmt.Errorf("failed to parse url: %w", err)
	}

	requestURL := parsed.String()
	if !parsed.IsAbs() {
		requestURL = fmt.Sprintf("%s%s", c.baseURL, path)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}

	return json.NewDecoder(resp.Body).Decode(response)
}
