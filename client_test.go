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

	l := pokeapi.List{
		Results: []pokeapi.ListResult{
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
		},
	}

	client, err := pokeapi.NewClient(srv.URL)
	require.NoError(t, err)

	res, err := pokeapi.FetchResults[pokeapi.Pokemon](t.Context(), client, &l)
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

	l := pokeapi.List{
		Results: []pokeapi.ListResult{
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
		},
	}

	client, err := pokeapi.NewClient(srv.URL)
	require.NoError(t, err)

	res, err := pokeapi.FetchResultsN[pokeapi.Pokemon](t.Context(), client, &l, 4)
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
