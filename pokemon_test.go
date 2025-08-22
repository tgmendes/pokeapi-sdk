package pokeapi_test

import (
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

	sdk := pokeapi.NewClient(srv.URL)
	pokemon, err := sdk.PokemonByName(t.Context(), "clefairy")
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

	sdk := pokeapi.NewClient(srv.URL)
	pokemon, err := sdk.PokemonByID(t.Context(), 35)
	require.NoError(t, err)
	require.NotNil(t, pokemon)
	assert.Equal(t, calledPath, "/pokemon/35")
	assert.Equal(t, calledMethod, http.MethodGet)
	assert.Equal(t, pokemon.ID, 35)
	assert.Equal(t, pokemon.Name, "clefairy")
}
