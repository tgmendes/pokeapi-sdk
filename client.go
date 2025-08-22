package pokeapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
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

func FetchResults[T any](ctx context.Context, c *Client, l *List) ([]T, error) {
	results := make([]T, 0, len(l.Results))
	for _, result := range l.Results {
		var resp T
		err := c.get(ctx, result.URL, &resp)
		if err != nil {
			return nil, err
		}
		results = append(results, resp)
	}

	return results, nil
}

func FetchResultsN[T any](ctx context.Context, c *Client, l *List, n int) ([]T, error) {
	if n < 1 {
		n = 1
	}
	type job struct {
		index int
		url   string
	}
	type res struct {
		index int
		value T
		err   error
	}

	jobs := make(chan job)
	out := make(chan res)

	var wg sync.WaitGroup
	for range n {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				var v T
				err := c.get(ctx, j.url, &v)
				out <- res{j.index, v, err}
			}
		}()
	}
	go func() {
		wg.Wait()
		close(out)
	}()

	go func() {
		for i, r := range l.Results {
			select {
			case <-ctx.Done():
				close(jobs)
				return
			case jobs <- job{i, r.URL}:
			}
		}
		close(jobs)
	}()

	results := make([]T, len(l.Results))
	for r := range out {
		if r.err != nil {
			return nil, r.err
		}
		results[r.index] = r.value
	}

	return results, ctx.Err()
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
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return HTTPError{
			Message:    string(b),
			StatusCode: resp.StatusCode,
		}
	}

	return json.NewDecoder(resp.Body).Decode(response)
}
