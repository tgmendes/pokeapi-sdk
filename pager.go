package pokeapi

import (
	"context"
	"errors"
)

var ErrNoMorePages = errors.New("no more pages")

type List struct {
	Count   int          `json:"count"`
	Next    *string      `json:"next"`
	Prev    *string      `json:"previous"`
	Results []ListResult `json:"results"`
}

type ListResult struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Pager is used to iterate over a list of results from the PokeAPI.
// This is not concurrent safe, and should be used sequentially.
type Pager struct {
	c        *Client
	next     *string
	previous *string
}

func NewPager(c *Client, startPath string) *Pager {
	return &Pager{
		c:    c,
		next: &startPath,
	}
}

func (p *Pager) Next(ctx context.Context) (*List, error) {
	if p.next == nil {
		return nil, ErrNoMorePages
	}

	return p.iter(ctx, *p.next)
}

func (p *Pager) Previous(ctx context.Context) (*List, error) {
	if p.previous == nil {
		return nil, ErrNoMorePages
	}

	return p.iter(ctx, *p.previous)
}

func (p *Pager) iter(ctx context.Context, url string) (*List, error) {
	var results List
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

	return &results, nil
}
