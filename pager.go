package pokeapi

import (
	"context"
	"errors"
)

var ErrNoMorePages = errors.New("no more pages")

type List[T any] struct {
	Count   int          `json:"count"`
	Next    *string      `json:"next"`
	Prev    *string      `json:"previous"`
	Results []ListResult `json:"results"`
}

type ListResult struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func (l List[T]) FetchResults(ctx context.Context, c *Client) ([]T, error) {
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

type Pager[T any] struct {
	c    *Client
	next *string
	prev *string
}

func NewPager[T any](c *Client, startPath string) *Pager[T] {
	return &Pager[T]{
		c:    c,
		next: &startPath,
	}
}

func (p *Pager[T]) Next(ctx context.Context) (*List[T], error) {
	if p.next == nil {
		return nil, ErrNoMorePages
	}

	return p.iter(ctx, *p.next)
}

func (p *Pager[T]) Previous(ctx context.Context) (*List[T], error) {
	if p.prev == nil {
		return nil, ErrNoMorePages
	}

	return p.iter(ctx, *p.prev)
}

func (p *Pager[T]) iter(ctx context.Context, url string) (*List[T], error) {
	var results List[T]
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
