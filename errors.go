package pokeapi

import "fmt"

type HTTPError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("http %s: %d", e.Message, e.StatusCode)
}
