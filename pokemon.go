package pokeapi

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/tgmendes/pokeapi-sdk/internal/apitype"
)

const pokemonPath = "/pokemon"

// Pokemon represents a Pokémon with its basic attributes and moves.
type Pokemon struct {
	ID             int    // Unique identifier for the Pokémon
	Name           string // Name of the Pokémon
	BaseExperience int    // Base experience gained when defeating this Pokémon
	Height         int    // Height in decimeters
	Order          int    // Order for sorting in the Pokédex
	Weight         int    // Weight
	IsLegendary    bool   // Whether this Pokémon is legendary
	IsMythical     bool   // Whether this Pokémon is mythical
	Moves          []Move // List of moves the Pokémon can learn
}

// Move represents a Pokémon move with its combat attributes.
type Move struct {
	Accuracy int    // Accuracy percentage of the move
	Name     string // Name of the move
	Power    int    // Base power of the move
	PP       int    // Power Points - how many times the move can be used
}

// PokemonByID fetches a Pokémon by its numeric ID.
func (c *Client) PokemonByID(ctx context.Context, id int) (*Pokemon, error) {
	path, err := url.JoinPath(pokemonPath, fmt.Sprint(id))
	if err != nil {
		return nil, fmt.Errorf("failed to join path: %w", err)
	}

	var apiPokemon apitype.Pokemon
	err = c.Get(ctx, path, &apiPokemon)
	if err != nil {
		return nil, fmt.Errorf("failed to get pokemon with id %d: %w", id, err)
	}

	return c.mapPokemon(ctx, apiPokemon)
}

// PokemonByName fetches a Pokémon by its name.
func (c *Client) PokemonByName(ctx context.Context, name string) (*Pokemon, error) {
	path, err := url.JoinPath(pokemonPath, name)
	if err != nil {
		return nil, fmt.Errorf("failed to join path: %w", err)
	}

	var pokemon apitype.Pokemon
	err = c.Get(ctx, path, &pokemon)
	if err != nil {
		return nil, fmt.Errorf("failed to get pokemon with name %s: %w", name, err)
	}

	return c.mapPokemon(ctx, pokemon)
}

// AllPokemon fetches all Pokémon using pagination.
// This method automatically handles pagination and fetches complete Pokémon data.
// Use with caution as it fetches a large amount of data.
func (c *Client) AllPokemon(ctx context.Context, options ...RequestOption) ([]Pokemon, error) {
	pager := c.PokemonPager(options...)

	var results []Pokemon
	for {
		list, err := pager.Next(ctx)
		if errors.Is(err, ErrNoMorePages) {
			break
		}

		if err != nil {
			return nil, err
		}

		apiPokemons, err := FetchResultsN[apitype.Pokemon](ctx, c, list, 5)
		if err != nil {
			return nil, err
		}

		pokemons := make([]Pokemon, 0, len(apiPokemons))
		for _, apiPokemon := range apiPokemons {
			pokemon, err := c.mapPokemon(ctx, apiPokemon)
			if err != nil {
				return nil, err
			}
			pokemons = append(pokemons, *pokemon)
		}

		results = append(results, pokemons...)
	}

	return results, nil
}

// PokemonPage fetches a single page of Pokémon resources.
// Returns a list of Resource objects containing names and URLs.
// You can parse these results to fetch pokemon data by using the FetchResults(N) function.
func (c *Client) PokemonPage(ctx context.Context, options ...RequestOption) ([]Resource, error) {
	pager := c.PokemonPager(options...)

	return pager.Next(ctx)
}

// PokemonPager creates a new pager for iterating through Pokémon pages.
// Use this for manual pagination control.
// You can parse these results to fetch pokemon data by using the FetchResults(N) function.
func (c *Client) PokemonPager(options ...RequestOption) *Pager {
	opts := defaultRequestOptions()
	if options != nil {
		opts = processOptions(options...)
	}

	q := opts.urlParams.Encode()
	startPath := pokemonPath
	if q != "" {
		startPath = fmt.Sprintf("%s?%s", pokemonPath, q)
	}

	return NewPager(c, startPath)
}

func (c *Client) mapPokemon(ctx context.Context, poke apitype.Pokemon) (*Pokemon, error) {
	movesResource := make([]Resource, 0, len(poke.Moves))
	for _, move := range poke.Moves {
		movesResource = append(movesResource, Resource{
			Name: move.Move.Name,
			URL:  move.Move.URL,
		})
	}

	moves, err := FetchResultsN[apitype.Move](ctx, c, movesResource, 5)
	if err != nil {
		return nil, err
	}

	speciesResource := Resource{
		Name: poke.Species.Name,
		URL:  poke.Species.URL,
	}

	species, err := FetchResults[apitype.Species](ctx, c, []Resource{speciesResource})
	if err != nil {
		return nil, err
	}

	pokeMoves := make([]Move, 0, len(moves))
	for _, move := range moves {
		pokeMoves = append(pokeMoves, Move{
			Accuracy: move.Accuracy,
			Name:     move.Name,
			Power:    move.Power,
			PP:       move.Pp,
		})
	}

	return &Pokemon{
		ID:             poke.ID,
		Name:           poke.Name,
		BaseExperience: poke.BaseExperience,
		Height:         poke.Height,
		Order:          poke.Order,
		Weight:         poke.Weight,
		IsLegendary:    species[0].IsLegendary,
		IsMythical:     species[0].IsMythical,
		Moves:          pokeMoves,
	}, nil
}
