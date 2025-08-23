package pokeapi

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/tgmendes/pokeapi-sdk/internal/apitype"
)

const generationPath = "/generation"

// Generation represents a Pokémon generation and some of its main features.
type Generation struct {
	ID             int      // Unique identifier for the generation
	Name           string   // Name of the generation (e.g., "generation-i")
	Region         string   // Main region of the generation
	Locations      []string // List of location names in this generation
	Moves          []string // List of move names introduced in this generation
	PokemonSpecies []string // List of Pokémon species introduced in this generation
	Types          []string // List of types introduced in this generation
}

// GenerationByID fetches a generation by its numeric ID.
func (c *Client) GenerationByID(ctx context.Context, id int) (*Generation, error) {
	path, err := url.JoinPath(generationPath, fmt.Sprint(id))
	if err != nil {
		return nil, fmt.Errorf("failed to join path: %w", err)
	}

	var apiGeneration apitype.Generation
	err = c.Get(ctx, path, &apiGeneration)
	if err != nil {
		return nil, fmt.Errorf("failed to get generation with id %d: %w", id, err)
	}

	return c.mapGeneration(ctx, apiGeneration)
}

// GenerationByName fetches a generation by its name (e.g., "generation-i").
func (c *Client) GenerationByName(ctx context.Context, name string) (*Generation, error) {
	path, err := url.JoinPath(generationPath, name)
	if err != nil {
		return nil, fmt.Errorf("failed to join path: %w", err)
	}

	var apiGeneration apitype.Generation
	err = c.Get(ctx, path, &apiGeneration)
	if err != nil {
		return nil, fmt.Errorf("failed to get generation with name %s: %w", name, err)
	}

	return c.mapGeneration(ctx, apiGeneration)
}

// AllGenerations fetches all generations using pagination.
// This method automatically handles pagination and fetches complete generation data.
func (c *Client) AllGenerations(ctx context.Context, options ...RequestOption) ([]Generation, error) {
	pager := c.GenerationPager(options...)

	var results []Generation
	for {
		list, err := pager.Next(ctx)
		if errors.Is(err, ErrNoMorePages) {
			break
		}

		if err != nil {
			return nil, err
		}

		apiGenerations, err := FetchResultsN[apitype.Generation](ctx, c, list, 5)
		if err != nil {
			return nil, err
		}

		generations := make([]Generation, 0, len(apiGenerations))
		for _, apiGeneration := range apiGenerations {
			generation, err := c.mapGeneration(ctx, apiGeneration)
			if err != nil {
				return nil, err
			}
			generations = append(generations, *generation)
		}

		results = append(results, generations...)
	}

	return results, nil
}

// GenerationPage fetches a single page of generation resources.
// Returns a list of Resource objects containing names and URLs.
// You can parse these results to fetch generation data by using the FetchResults(N) function.
func (c *Client) GenerationPage(ctx context.Context, options ...RequestOption) ([]Resource, error) {
	pager := c.GenerationPager(options...)

	return pager.Next(ctx)
}

// GenerationPager creates a new pager for iterating through generation pages.
// Use this for manual pagination control.
// You can parse these results to fetch generation data by using the FetchResults(N) function.
func (c *Client) GenerationPager(options ...RequestOption) *Pager {
	opts := defaultRequestOptions()
	if options != nil {
		opts = processOptions(options...)
	}

	q := opts.urlParams.Encode()
	startPath := generationPath
	if q != "" {
		startPath = fmt.Sprintf("%s?%s", generationPath, q)
	}

	return NewPager(c, startPath)
}

func (c *Client) mapGeneration(ctx context.Context, gen apitype.Generation) (*Generation, error) {
	moves := make([]string, 0, len(gen.Moves))
	for _, move := range gen.Moves {
		moves = append(moves, move.Name)
	}

	species := make([]string, 0, len(gen.PokemonSpecies))
	for _, apiSpecies := range gen.PokemonSpecies {
		species = append(species, apiSpecies.Name)
	}

	pokeTypes := make([]string, 0, len(gen.Types))
	for _, apiType := range gen.Types {
		pokeTypes = append(pokeTypes, apiType.Name)
	}

	regionResource := Resource{
		Name: gen.MainRegion.Name,
		URL:  gen.MainRegion.URL,
	}

	region, err := FetchResults[apitype.Region](ctx, c, []Resource{regionResource})
	if err != nil {
		return nil, err
	}

	locationNames := make([]string, 0, len(region[0].Locations))
	for _, location := range region[0].Locations {
		locationNames = append(locationNames, location.Name)
	}

	return &Generation{
		ID:             gen.ID,
		Name:           gen.Name,
		Region:         gen.MainRegion.Name,
		Locations:      locationNames,
		Moves:          moves,
		PokemonSpecies: species,
		Types:          pokeTypes,
	}, nil
}
