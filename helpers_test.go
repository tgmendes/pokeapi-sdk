package pokeapi_test

import (
	"embed"
	"testing"
)

//go:embed fixtures
var fixtures embed.FS

func fixture(t *testing.T, name string) []byte {
	data, err := fixtures.ReadFile("fixtures/" + name)
	if err != nil {
		t.Fatalf("failed to read fixture %s: %v", name, err)
	}
	return data
}
