package pokeapi

import (
	"net/url"
	"strconv"
)

const (
	defaultLimit  = 20
	defaultOffset = 0
)

// RequestOption is a function that configures request parameters.
type RequestOption func(*requestOptions)

type requestOptions struct {
	urlParams url.Values
}

func defaultRequestOptions() requestOptions {
	return processOptions(Limit(20), Offset(0))
}

// Limit sets the number of entries that a request should return.
// The default limit is 20 entries per page.
func Limit(amount int) RequestOption {
	return func(o *requestOptions) {
		o.urlParams.Set("limit", strconv.Itoa(amount))
	}
}

// Offset sets the index of the first entry to return.
// Use this for pagination to skip a certain number of entries.
func Offset(amount int) RequestOption {
	return func(o *requestOptions) {
		o.urlParams.Set("offset", strconv.Itoa(amount))
	}
}

func processOptions(options ...RequestOption) requestOptions {
	o := requestOptions{
		urlParams: url.Values{},
	}
	for _, opt := range options {
		opt(&o)
	}

	return o
}
