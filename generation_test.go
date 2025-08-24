package pokeapi_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tgmendes/pokeapi-sdk"
)

func TestGetGenerationByName(t *testing.T) {
	srv := setupGenerationServer(t)
	defer srv.Close()

	client, err := pokeapi.NewClient(srv.URL)
	require.NoError(t, err)

	generation, err := client.GenerationByName(t.Context(), "generation-i")
	require.NoError(t, err)
	require.NotNil(t, generation)
	assert.Equal(t, 1, generation.ID)
	assert.Equal(t, "generation-i", generation.Name)
	assert.NotEmpty(t, generation.Moves)
	assert.Contains(t, generation.Moves, "confusion")
	assert.Contains(t, generation.PokemonSpecies, "bulbasaur")
	assert.Contains(t, generation.Types, "fire")
	assert.Contains(t, generation.Locations, "celadon-city")
}

func TestGetGenerationByID(t *testing.T) {
	srv := setupGenerationServer(t)
	defer srv.Close()

	client, err := pokeapi.NewClient(srv.URL)
	require.NoError(t, err)

	generation, err := client.GenerationByID(t.Context(), 1)
	require.NoError(t, err)
	require.NotNil(t, generation)
	assert.Equal(t, 1, generation.ID)
	assert.Equal(t, "generation-i", generation.Name)
}

func TestGenerationPage(t *testing.T) {
	srv := setupGenerationServer(t)
	defer srv.Close()

	client, err := pokeapi.NewClient(srv.URL)
	require.NoError(t, err)

	results, err := client.GenerationPage(t.Context())
	require.NoError(t, err)
	require.NotNil(t, results)
	assert.Len(t, results, 9)
	assert.Equal(t, "generation-i", results[0].Name)
}

func TestAllGeneration(t *testing.T) {
	srv := setupGenerationServer(t)
	defer srv.Close()

	client, err := pokeapi.NewClient(srv.URL)
	require.NoError(t, err)

	generation, err := client.AllGenerations(t.Context())
	require.NoError(t, err)
	require.NotNil(t, generation)
	assert.Len(t, generation, 9)
	assert.Equal(t, 1, generation[0].ID)
	assert.Equal(t, "generation-i", generation[0].Name)
}

func TestGenerationPager(t *testing.T) {
	gotLimit := ""
	gotOffset := ""
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			gotLimit = r.URL.Query().Get("limit")
			gotOffset = r.URL.Query().Get("offset")
			_, _ = w.Write(fixture(t, "generation_page1.json"))
		},
	))

	tests := []struct {
		name      string
		opts      []pokeapi.RequestOption
		expLimit  string
		expOffset string
	}{
		{
			name:      "uses default options when no options are passed",
			expLimit:  "20",
			expOffset: "0",
		},
		{
			name: "uses limit when limit is passed",
			opts: []pokeapi.RequestOption{
				pokeapi.Limit(5),
			},
			expLimit:  "5",
			expOffset: "",
		},
		{
			name: "uses limit and offset when both are passed",
			opts: []pokeapi.RequestOption{
				pokeapi.Limit(5),
				pokeapi.Offset(5),
			},
			expLimit:  "5",
			expOffset: "5",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotLimit = ""
			gotOffset = ""
			client, err := pokeapi.NewClient(srv.URL)
			require.NoError(t, err)
			pager := client.GenerationPager(test.opts...)
			_, err = pager.Next(t.Context())
			require.NoError(t, err)
			assert.Equal(t, test.expLimit, gotLimit)
			assert.Equal(t, test.expOffset, gotOffset)
		})
	}
}

func TestParseGenerationResource(t *testing.T) {
	srv := setupGenerationServer(t)
	defer srv.Close()

	client, err := pokeapi.NewClient(srv.URL)
	require.NoError(t, err)

	page, err := client.GenerationPage(t.Context())
	require.NoError(t, err)

	generations, err := client.ParseGenerationResource(t.Context(), page)
	require.NoError(t, err)
	require.NotNil(t, generations)
	assert.Len(t, generations, 9)
}

func setupGenerationServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			switch {
			case strings.Contains(r.URL.Path, "region"):
				_, _ = w.Write(fixture(t, "region.json"))
				return
			case r.URL.Query().Get("offset") == "0":
				_, _ = w.Write(fixture(t, "generation_page1.json"))
				return
			default:
				_, _ = w.Write(fixture(t, "generation.json"))
			}
		},
	))
}
