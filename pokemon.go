package pokeapi

import (
	"context"
	"errors"
	"fmt"

	"github.com/tgmendes/pokeapi-sdk/internal/apitype"
)

type Pokemon struct {
	ID             int
	Name           string
	BaseExperience int
	Height         int
	Order          int
	Weight         int
	IsLegendary    bool
	IsMythical     bool
	Moves          []Move
}

type Move struct {
	Accuracy int
	Name     string
	Power    int
	PP       int
}

const pokemonPath = "/pokemon"

func (c *Client) PokemonByID(ctx context.Context, id int) (*Pokemon, error) {
	path := fmt.Sprintf("%s/%d", pokemonPath, id)
	var apiPokemon apitype.Pokemon
	err := c.Get(ctx, path, &apiPokemon)
	if err != nil {
		return nil, fmt.Errorf("failed to get pokemon with id %d: %w", id, err)
	}

	return c.mapPokemon(ctx, apiPokemon)
}

func (c *Client) PokemonByName(ctx context.Context, name string) (*Pokemon, error) {
	path := fmt.Sprintf("%s/%s", pokemonPath, name)

	var pokemon apitype.Pokemon
	err := c.Get(ctx, path, &pokemon)
	if err != nil {
		return nil, fmt.Errorf("failed to get pokemon with name %s: %w", name, err)
	}

	return c.mapPokemon(ctx, pokemon)
}

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

func (c *Client) PokemonPage(ctx context.Context, options ...RequestOption) ([]Resource, error) {
	pager := c.PokemonPager(options...)

	return pager.Next(ctx)
}

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
