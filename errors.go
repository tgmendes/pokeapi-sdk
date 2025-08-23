package pokeapi

import "fmt"

// HTTPError represents an HTTP error response from the API.
type HTTPError struct {
	Message    string `json:"message"`     // Error message from the API
	StatusCode int    `json:"status_code"` // HTTP status code
}

// Error implements the error interface for HTTPError.
func (e HTTPError) Error() string {
	return fmt.Sprintf("http %s: %d", e.Message, e.StatusCode)
}
