//go:build integration

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
	assert.Equal(t, 1, pokemon.ID)
	assert.Equal(t, "bulbasaur", pokemon.Name)
}

func TestPokemonByName_Integration(t *testing.T) {
	client, err := pokeapi.NewClient("https://pokeapi.co/api/v2")

	pokemon, err := client.PokemonByName(t.Context(), "bulbasaur")
	require.NoError(t, err)
	assert.Equal(t, 1, pokemon.ID)
	assert.Equal(t, "bulbasaur", pokemon.Name)
}

func TestAllPokemon_Integration(t *testing.T) {
	client, err := pokeapi.NewClient(
		"https://pokeapi.co/api/v2",
		pokeapi.WithLimit(2000, 2000))

	pokemon, err := client.AllPokemon(t.Context(), pokeapi.Limit(100))
	require.NoError(t, err)
	assert.Len(t, pokemon, 1302,
		"should have 1302 pokemon but got %d", len(pokemon))
	assert.Equal(t, 1, pokemon[0].ID)
	assert.Equal(t, "bulbasaur", pokemon[0].Name)
}

func TestGenerationByID_Integration(t *testing.T) {
	client, err := pokeapi.NewClient("https://pokeapi.co/api/v2")

	generation, err := client.GenerationByID(t.Context(), 1)
	require.NoError(t, err)
	assert.Equal(t, 1, generation.ID)
	assert.Equal(t, "generation-i", generation.Name)
}

func TestGenerationByName_Integration(t *testing.T) {
	client, err := pokeapi.NewClient("https://pokeapi.co/api/v2")

	generation, err := client.GenerationByName(t.Context(), "generation-i")
	require.NoError(t, err)
	assert.Equal(t, 1, generation.ID)
	assert.Equal(t, "generation-i", generation.Name)
}

func TestAllGeneration_Integration(t *testing.T) {
	client, err := pokeapi.NewClient("https://pokeapi.co/api/v2")

	generation, err := client.AllGenerations(t.Context())
	require.NoError(t, err)
	assert.Len(t, generation, 9)
	assert.Equal(t, 1, generation[0].ID)
	assert.Equal(t, "generation-i", generation[0].Name)
}
