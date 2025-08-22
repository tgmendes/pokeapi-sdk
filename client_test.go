package pokeapi_test

import (
	"encoding/json"
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

	c := pokeapi.NewClient(srv.URL)
	res, err := pokeapi.FetchListResults[pokeapi.Pokemon](t.Context(), c, &l)
	require.NoError(t, err)

	assert.Len(t, res, 3)
}
