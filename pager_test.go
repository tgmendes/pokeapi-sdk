package pokeapi_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tgmendes/pokeapi-sdk"
)

func TestPager_Next(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// because it's a test we can hardcode the comparison - if we wanted something more
			// robust we could parse the URL and compare the query parameters.
			if r.URL.Query().Get("offset") == "0" {
				_, _ = w.Write(fixture(t, "pokemon_page1.json"))
				return
			}

			if r.URL.Query().Get("offset") == "20" {
				_, _ = w.Write(fixture(t, "pokemon_page2.json"))
				return
			}
		},
	))
	defer srv.Close()

	sdk := pokeapi.NewClient(srv.URL)

	startPath := fmt.Sprintf("%s/pokemon?offset=0", srv.URL)
	pager := pokeapi.NewPager[pokeapi.Pokemon](sdk, startPath)

	var totalResults []pokeapi.List[pokeapi.Pokemon]
	for {
		list, err := pager.Next(t.Context())
		if errors.Is(err, pokeapi.ErrNoMorePages) {
			break
		}
		require.NoError(t, err)
		require.NotNil(t, list)
		totalResults = append(totalResults, *list)
	}

	assert.Len(t, totalResults, 2)
}

func TestPager_Previous(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// because it's a test we can hardcode the comparison - if we wanted something more
			// robust we could parse the URL and compare the query parameters.
			if r.URL.Query().Get("offset") == "0" {
				_, _ = w.Write(fixture(t, "pokemon_page1.json"))
				return
			}

			if r.URL.Query().Get("offset") == "20" {
				_, _ = w.Write(fixture(t, "pokemon_page2.json"))
				return
			}
		},
	))
	defer srv.Close()

	sdk := pokeapi.NewClient(srv.URL)

	startPath := fmt.Sprintf("%s/pokemon?offset=0", srv.URL)
	pager := pokeapi.NewPager[pokeapi.Pokemon](sdk, startPath)

	// fetch the first page
	_, err := pager.Next(t.Context())
	require.NoError(t, err)

	// there is no previous page before the first page
	_, err = pager.Previous(t.Context())
	require.ErrorIs(t, err, pokeapi.ErrNoMorePages)

	// fetch the second page
	_, err = pager.Next(t.Context())
	require.NoError(t, err)

	// now there is a previous page, and it is the first page
	list, err := pager.Previous(t.Context())
	require.NoError(t, err)
	require.NotNil(t, list)
	// first element of first page is bulbasaur
	assert.Equal(t, list.Results[0].Name, "bulbasaur")
}

func TestList_FetchResults(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			res := pokeapi.Pokemon{}

			err := json.NewEncoder(w).Encode(&res)
			require.NoError(t, err)
		},
	))
	defer srv.Close()

	l := pokeapi.List[pokeapi.Pokemon]{
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
	res, err := l.FetchResults(t.Context(), c)
	require.NoError(t, err)

	assert.Len(t, res, 3)
}
