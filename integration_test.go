package pokeapi_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tgmendes/pokeapi-sdk"
)

func TestPokemonByID_Integration(t *testing.T) {
	client, err := pokeapi.NewClient("https://pokeapi.co/api/v2")

	pokemon, err := client.PokemonByID(t.Context(), 1)
	require.NoError(t, err)
	assert.Equal(t, pokemon.ID, 1)
	assert.Equal(t, pokemon.Name, "bulbasaur")
}

func TestPokemonByName_Integration(t *testing.T) {
	client, err := pokeapi.NewClient("https://pokeapi.co/api/v2")

	pokemon, err := client.PokemonByName(t.Context(), "bulbasaur")
	require.NoError(t, err)
	assert.Equal(t, pokemon.ID, 1)
	assert.Equal(t, pokemon.Name, "bulbasaur")
}

func TestAllPokemon_Integration(t *testing.T) {
	client, err := pokeapi.NewClient(
		"https://pokeapi.co/api/v2",
		pokeapi.WithLimit(2000, 2000))

	pokemon, err := client.AllPokemon(t.Context(), pokeapi.Limit(100))
	require.NoError(t, err)
	assert.Len(t, pokemon, 1302,
		"should have 1302 pokemon but got %d", len(pokemon))
	assert.Equal(t, pokemon[0].ID, 1)
	assert.Equal(t, pokemon[0].Name, "bulbasaur")
}
