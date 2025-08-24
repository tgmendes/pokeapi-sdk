# PokeAPI SDK

A Go SDK for interacting with the [PokéAPI](https://pokeapi.co/). This library provides a simple, 
efficient way to fetch Pokemon data with built-in caching, rate limiting, and structured data types.

## Installation

```bash
go get github.com/tgmendes/pokeapi-sdk
```

## Features

This SDK allows any consumer to interact with the pokemon API.

For now it is limited in scope, and only provides methods to fetch pokemon
or generations. 

Results from the API are parsed, enriched and transformed to provide a more 
ergonomic experience (for example, pokemon moves data is automatically fetched
and transformed to a more readable format).

Summary of features:
- Get Pokemon or Generation data
- Fetch resources by name, id, or all
- Paginate through lists of results from the API
- Rate limiting with configurable burst and limit
- Automatic caching of results
- Access Pokemon moves, stats, and other attributes

## Usage

You find an example in the [cmd/example.go](cmd/example.go) file, and run it using:

```bash
go run cmd/example.go
```

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/tgmendes/pokeapi-sdk"
)

func main() {
	client, _ := pokeapi.NewClient("https://pokeapi.co/api/v2",
		pokeapi.WithLimit(1000, 2000))

	// get pokemon by name
	pokemon, _ := client.PokemonByName(context.Background(), "pikachu")

	fmt.Printf("Pokémon: %s (ID: %d)\n", pokemon.Name, pokemon.ID)

	// get pokemon by id
	pokemon, _ = client.PokemonByID(context.Background(), 1)

	fmt.Printf("Pokémon: %s (ID: %d)\n", pokemon.Name, pokemon.ID)

	// get all pokemon
	allPokemon, _ := client.AllPokemon(context.Background())

	fmt.Printf("First Pokémon: %s (ID: %d)\n", allPokemon[0].Name, allPokemon[0].ID)
	fmt.Printf("Last Pokémon: %s (ID: %d)\n",
		allPokemon[len(allPokemon)-1].Name, allPokemon[len(allPokemon)-1].ID)

	// get generation by id
	generation, _ := client.GenerationByID(context.Background(), 1)

	fmt.Printf("Generation: %s (ID: %d)\n", generation.Name, generation.ID)

	// get generation by name
	generation, _ = client.GenerationByName(context.Background(), "generation-i")

	fmt.Printf("Generation: %s (ID: %d)\n", generation.Name, generation.ID)

	// get all generations
	generations, _ := client.AllGenerations(context.Background())

	fmt.Printf("First Generation: %s (ID: %d)\n", generations[0].Name, generations[0].ID)

	// paginating
	pager := client.PokemonPager(pokeapi.Limit(50))
	firstPage, _ := pager.Next(context.Background())
	res, _ := pokeapi.FetchResultsN[pokeapi.Pokemon](context.Background(), client, firstPage, 5)
	fmt.Printf("First Pokemon: %d\n", len(res[0].Name))

	// next page
	secondPage, _ := pager.Next(context.Background())
	res, _ = pokeapi.FetchResultsN[pokeapi.Pokemon](context.Background(), client, secondPage, 5)
	fmt.Printf("Second Pokemon: %d\n", len(res[0].Name))

	// previous page
	previousPage, _ := pager.Previous(context.Background())
	res, _ = pokeapi.FetchResultsN[pokeapi.Pokemon](context.Background(), client, previousPage, 5)
	fmt.Printf("Previous Pokemon: %d\n", len(res[0].Name))    
}
```

## Testing

To run unit tests:

```bash
go test -v ./...
```

To run integration tests:

```bash
go test -v -tags=integration ./...
```

## Tools used

- [Go 1.24](https://golang.org)
- [PokeAPI](https://pokeapi.co) as the main data source
- [Testify](https://pkg.go.dev/github.com/stretchr/testify) for assertions
- [QucikType](https://app.quicktype.io/) for generating API structs
- [Ristretto](https://github.com/dgraph-io/ristretto) for caching

## Design Philosophy

The SDK was designed to be simple but extensible. For a client of this type
a flat structure was chosen, with just a few internals (cache and api types). This
keeps things simple, while also being testable.

The API types being in a non-importable package allows us to cleary define 
the boundaries between what is the raw data types from the API, and what is the 
public objects we want to expose to consumers (which are not just the plain API response, but
are enriched to provide a better user experience).

When multiple requests are involved (such as fetching list resource), concurrency
was chosen to avoid requests taking too long. Our rate limiting logic should ensure that we are not 
overloading the API.

## Future Improvements

* Tweak rate limiting
* Improved error handling and return codes
* Improved caching strategy
* More comprehensive tests
* Automatic generation of methods and API structs
* Better handling of query parameters (use URL) and setup parameters on generic 
client rather than on each method
* Create a test server that can handle all requests
* Add fun methods to pokemon (e.g. `pokemon.canBeat(otherPokemon)`)
* Configurable parameters for fetching results