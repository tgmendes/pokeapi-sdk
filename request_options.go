package pokeapi

import (
	"net/url"
	"strconv"
)

const (
	defaultLimit  = 20
	defaultOffset = 0
)

type RequestOption func(*requestOptions)

type requestOptions struct {
	urlParams url.Values
}

func defaultRequestOptions() requestOptions {
	return processOptions(Limit(20), Offset(0))
}

// Limit sets the number of entries that a request should return.
func Limit(amount int) RequestOption {
	return func(o *requestOptions) {
		o.urlParams.Set("limit", strconv.Itoa(amount))
	}
}

// Offset sets the index of the first entry to return.
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
