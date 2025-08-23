// Package pokeapi provides a Go SDK for interacting with the PokéAPI (https://pokeapi.co/).
//
// This package offers a simple, efficient way to fetch Pokémon data with built-in
// caching, rate limiting, and structured data types. It handles pagination automatically
// and provides both sequential and concurrent data fetching options.
//
// Basic usage:
//
//	client, err := pokeapi.NewClient("https://pokeapi.co/api/v2", pokeapi.WithLimit(10, 20))
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	pokemon, err := client.PokemonByName(context.Background(), "pikachu")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Printf("Pokémon: %s (ID: %d)\n", pokemon.Name, pokemon.ID)
package pokeapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"github.com/tgmendes/pokeapi-sdk/internal/cache"
)

// Client represents a PokéAPI client with built-in caching and rate limiting.
type Client struct {
	http    *http.Client
	baseURL string
	cache   *cache.Cache
	limiter *rate.Limiter
}

// ClientOption is a function that configures a Client.
type ClientOption func(*Client)

// WithLimit is a client option that configures the rate limiter for the client.
// rpsLimit sets the requests per second limit, burstLimit sets the burst size.
func WithLimit(rpsLimit float64, burstLimit int) ClientOption {
	return func(c *Client) {
		c.limiter = rate.NewLimiter(rate.Limit(rpsLimit), burstLimit)
	}
}

// NewClient creates a new PokéAPI client with the given base URL and options.
// The client includes default rate limiting (20 RPS, 50 burst) and caching.
func NewClient(baseURL string, opts ...ClientOption) (*Client, error) {
	memCache, err := cache.New()
	if err != nil {
		return nil, err
	}

	c := Client{
		http: &http.Client{
			Timeout: 5 * time.Second,
		},
		baseURL: baseURL,
		cache:   memCache,
		limiter: rate.NewLimiter(rate.Limit(20), 50),
	}

	for _, opt := range opts {
		opt(&c)
	}

	return &c, nil
}

// Get performs a GET request to the given path and unmarshals the response.
// The path can be absolute or relative to the client's base URL. If a relative path is given,
// an absolute path is constructed by joining the base URL and the relative path.
// Responses are cached and rate limited automatically.
func (c *Client) Get(ctx context.Context, path string, response any) error {
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

	cached, ok := c.cache.Get(cacheKey(requestURL))
	if ok {
		return json.Unmarshal(cached, response)
	}

	if err = c.limiter.Wait(ctx); err != nil {
		return fmt.Errorf("failed to wait for rate limit: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent",
		"pokeapi-sdk/1.0 (+https://github.com/tgmendes/pokeapi-sdk)")

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

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	c.cache.Set(cacheKey(requestURL), b)

	return json.Unmarshal(b, response)
}

// FetchResults fetches all resources from the given list sequentially.
// It makes individual API calls for each resource URL and returns the results.
func FetchResults[T any](ctx context.Context, c *Client, l []Resource) ([]T, error) {
	results := make([]T, 0, len(l))
	for _, result := range l {
		var resp T
		err := c.Get(ctx, result.URL, &resp)
		if err != nil {
			return nil, err
		}
		results = append(results, resp)
	}

	return results, nil
}

// FetchResultsN fetches all resources from the given list concurrently using n workers.
// It's more efficient than FetchResults for large lists but uses more resources.
func FetchResultsN[T any](ctx context.Context, c *Client, l []Resource, n int) ([]T, error) {
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

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	jobs := make(chan job)
	out := make(chan res)

	var wg sync.WaitGroup
	for range n {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				var v T
				err := c.Get(ctx, j.url, &v)
				out <- res{j.index, v, err}
			}
		}()
	}
	go func() {
		wg.Wait()
		close(out)
	}()

	go func() {
		for i, r := range l {
			select {
			case <-ctx.Done():
				close(jobs)
				return
			case jobs <- job{i, r.URL}:
			}
		}
		close(jobs)
	}()

	results := make([]T, len(l))
	var firstErr error
	for r := range out {
		if r.err != nil && firstErr == nil {
			firstErr = r.err
			// Stop everyone else; keep draining 'out' until closed.
			cancel()
			continue
		}
		if firstErr == nil { // ignore successes after cancel
			results[r.index] = r.value
		}
	}

	if firstErr != nil {
		return nil, firstErr
	}
	return results, ctx.Err()
}

func cacheKey(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		// fallback: just return url
		return rawURL
	}

	// sort query params
	q := u.Query()
	keys := make([]string, 0, len(q))
	for k := range q {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	values := url.Values{}
	for _, k := range keys {
		vs := q[k]
		sort.Strings(vs)
		for _, v := range vs {
			values.Add(k, v)
		}
	}
	u.RawQuery = values.Encode()

	return u.String()
}
