package pokeapi_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tgmendes/pokeapi-sdk"
)

func TestGetPokemonByName(t *testing.T) {
	calledPath := ""
	calledMethod := ""
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			calledPath = r.URL.Path
			calledMethod = r.Method
			_, _ = w.Write(fixture(t, "clefairy.json"))
		},
	))
	defer srv.Close()

	client := pokeapi.NewClient(srv.URL)
	pokemon, err := client.PokemonByName(t.Context(), "clefairy")
	require.NoError(t, err)
	require.NotNil(t, pokemon)
	assert.Equal(t, calledPath, "/pokemon/clefairy")
	assert.Equal(t, calledMethod, http.MethodGet)
	assert.Equal(t, pokemon.ID, 35)
	assert.Equal(t, pokemon.Name, "clefairy")
}

func TestGetPokemonByID(t *testing.T) {
	calledPath := ""
	calledMethod := ""
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			calledPath = r.URL.Path
			calledMethod = r.Method
			_, _ = w.Write(fixture(t, "clefairy.json"))
		},
	))
	defer srv.Close()

	client := pokeapi.NewClient(srv.URL)
	pokemon, err := client.PokemonByID(t.Context(), 35)
	require.NoError(t, err)
	require.NotNil(t, pokemon)
	assert.Equal(t, calledPath, "/pokemon/35")
	assert.Equal(t, calledMethod, http.MethodGet)
	assert.Equal(t, pokemon.ID, 35)
	assert.Equal(t, pokemon.Name, "clefairy")
}

func TestPokemonPage(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write(fixture(t, "pokemon_page1.json"))
		},
	))
	defer srv.Close()

	client := pokeapi.NewClient(srv.URL)
	page, err := client.PokemonPage(t.Context())
	require.NoError(t, err)
	require.NotNil(t, page)
	assert.Len(t, page.Results, 20)
	assert.Equal(t, page.Results[0].Name, "bulbasaur")
}

func TestAllPokemon(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Println(r.URL.Query().Get("offset"))
			// because it's a test we can hardcode the comparison - if we wanted something more
			// robust we could parse the URL and compare the query parameters.
			if r.URL.Query().Get("offset") == "0" {
				fmt.Println("here")
				_, _ = w.Write(fixture(t, "pokemon_page1.json"))
				return
			}

			if r.URL.Query().Get("offset") == "20" {
				_, _ = w.Write(fixture(t, "pokemon_page2.json"))
				return
			}

			// this is a direct query to fetch a pokemon
			if r.URL.Query().Get("offset") == "" {
				// for testing purposes just write same pokemon all the time
				_, _ = w.Write(fixture(t, "clefairy.json"))
				return
			}
		},
	))
	defer srv.Close()

	client := pokeapi.NewClient(srv.URL)
	pokemon, err := client.AllPokemon(t.Context())
	require.NoError(t, err)
	require.NotNil(t, pokemon)
	assert.Len(t, pokemon, 40)
	assert.Equal(t, pokemon[0].ID, 35)
	assert.Equal(t, pokemon[0].Name, "clefairy")
}

func TestPokemonPager(t *testing.T) {
	gotLimit := ""
	gotOffset := ""
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			gotLimit = r.URL.Query().Get("limit")
			gotOffset = r.URL.Query().Get("offset")
			_, _ = w.Write(fixture(t, "pokemon_page1.json"))
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
			client := pokeapi.NewClient(srv.URL)
			pager := client.PokemonPager(test.opts...)
			_, err := pager.Next(t.Context())
			require.NoError(t, err)
			assert.Equal(t, test.expLimit, gotLimit)
			assert.Equal(t, test.expOffset, gotOffset)
		})
	}
}
