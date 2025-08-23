package pokeapi

import (
	"context"
	"errors"

	"github.com/tgmendes/pokeapi-sdk/internal/apitype"
)

// Resource represents a basic resource with a name and URL.
// It maps to an API named resource and is used to fetch the full resource data.
type Resource struct {
	Name string // Name of the resource
	URL  string // URL to fetch the full resource data
}

// ErrNoMorePages is returned when there are no more pages to iterate through.
var ErrNoMorePages = errors.New("no more pages")

// Pager is used to iterate over a list of results from the PokeAPI.
// This is not concurrent safe, and should be used sequentially.
type Pager struct {
	c        *Client
	next     *string
	previous *string
}

// NewPager creates a new pager starting from the given path.
// The pager can be used to iterate through paginated API results.
func NewPager(c *Client, startPath string) *Pager {
	return &Pager{
		c:    c,
		next: &startPath,
	}
}

// Next fetches the next page of results.
// Returns ErrNoMorePages when there are no more pages available.
func (p *Pager) Next(ctx context.Context) ([]Resource, error) {
	if p.next == nil {
		return nil, ErrNoMorePages
	}

	return p.iter(ctx, *p.next)
}

// Previous fetches the previous page of results.
// Returns ErrNoMorePages when there are no more pages available.
func (p *Pager) Previous(ctx context.Context) ([]Resource, error) {
	if p.previous == nil {
		return nil, ErrNoMorePages
	}

	return p.iter(ctx, *p.previous)
}

func (p *Pager) iter(ctx context.Context, url string) ([]Resource, error) {
	var results apitype.List
	if err := p.c.Get(ctx, url, &results); err != nil {
		return nil, err
	}

	if results.Next == nil || *results.Next == "" {
		p.next = nil
	} else {
		// avoid using same pointer as the list results
		next := *results.Next
		p.next = &next
	}

	if results.Prev == nil || *results.Prev == "" {
		p.previous = nil
	} else {
		prev := *results.Prev
		p.previous = &prev
	}

	resource := make([]Resource, 0, len(results.Results))
	for _, result := range results.Results {
		resource = append(resource, Resource{
			Name: result.Name,
			URL:  result.URL,
		})
	}
	return resource, nil
}
