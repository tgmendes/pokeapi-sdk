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

type Pager struct {
	c    *Client
	next *string
	prev *string
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
	if p.prev == nil {
		return nil, ErrNoMorePages
	}

	return p.iter(ctx, *p.prev)
}

func (p *Pager) iter(ctx context.Context, url string) (*List, error) {
	var results List
	if err := p.c.get(ctx, url, &results); err != nil {
		return nil, err
	}

	if results.Next == nil || *results.Next == "" {
		p.next = nil
	} else {
		p.next = results.Next
	}

	if results.Prev == nil || *results.Prev == "" {
		p.prev = nil
	} else {
		p.prev = results.Prev
	}

	return &results, nil
}
