package pokeapi_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tgmendes/pokeapi-sdk"
)

func TestList_FetchResults(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			res := pokeapi.Pokemon{}

			err := json.NewEncoder(w).Encode(&res)
			require.NoError(t, err)
		},
	))
	defer srv.Close()

	l := []pokeapi.Resource{
		{
			Name: "bulbasaur",
			URL:  "/pokemon/1",
		},
		{
			Name: "ivysaur",
			URL:  "/pokemon/1",
		},
		{
			Name: "venusaur",
			URL:  "/pokemon/1",
		},
	}

	client, err := pokeapi.NewClient(srv.URL)
	require.NoError(t, err)

	res, err := pokeapi.FetchResults[pokeapi.Pokemon](t.Context(), client, l)
	require.NoError(t, err)

	assert.Len(t, res, 3)
}

func TestList_FetchResultsN(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			res := pokeapi.Pokemon{}

			err := json.NewEncoder(w).Encode(&res)
			require.NoError(t, err)
		},
	))
	defer srv.Close()

	l := []pokeapi.Resource{
		{
			Name: "bulbasaur",
			URL:  "/pokemon/1",
		},
		{
			Name: "ivysaur",
			URL:  "/pokemon/1",
		},
		{
			Name: "venusaur",
			URL:  "/pokemon/1",
		},
		{
			Name: "ivysaur",
			URL:  "/pokemon/1",
		},
		{
			Name: "venusaur",
			URL:  "/pokemon/1",
		},
		{
			Name: "ivysaur",
			URL:  "/pokemon/1",
		},
		{
			Name: "venusaur",
			URL:  "/pokemon/1",
		},
		{
			Name: "ivysaur",
			URL:  "/pokemon/1",
		},
		{
			Name: "venusaur",
			URL:  "/pokemon/1",
		},
		{
			Name: "ivysaur",
			URL:  "/pokemon/1",
		},
		{
			Name: "venusaur",
			URL:  "/pokemon/1",
		},
	}

	client, err := pokeapi.NewClient(srv.URL)
	require.NoError(t, err)

	res, err := pokeapi.FetchResultsN[pokeapi.Pokemon](t.Context(), client, l, 4)
	require.NoError(t, err)

	assert.Len(t, res, 11)
}

func TestClient_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("something went wrong"))
		},
	))
	defer srv.Close()

	client, err := pokeapi.NewClient(srv.URL)
	require.NoError(t, err)

	_, err = client.PokemonPage(t.Context())
	require.Error(t, err)

	var httpErr pokeapi.HTTPError
	errors.As(err, &httpErr)
	assert.Equal(t, httpErr.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, httpErr.Message, "something went wrong")
}

func TestClient_Cache(t *testing.T) {
	serverCalled := 0
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			serverCalled++
			_, _ = w.Write(fixture(t, "pokemon_page1.json"))
		},
	))
	defer srv.Close()

	client, err := pokeapi.NewClient(srv.URL)
	require.NoError(t, err)

	reqList, err := client.PokemonPage(t.Context())
	require.NoError(t, err)
	require.Equal(t, serverCalled, 1)

	// fetch the same page again and make sure it doesn't call the server
	cachedList, err := client.PokemonPage(t.Context())
	require.NoError(t, err)
	assert.Equal(t, 1, serverCalled)
	assert.Equal(t, cachedList, reqList)
}
