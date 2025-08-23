package main

import (
	"context"
	"fmt"

	"github.com/tgmendes/pokeapi-sdk"
)

func main() {
	client, _ := pokeapi.NewClient("https://pokeapi.co/api/v2",
		pokeapi.WithLimit(1000, 2000))

	// get pokemon by name
	pokemon, _ := client.PokemonByName(context.Background(), "pikachu")

	fmt.Printf("Pokémon: %s (ID: %d)\n", pokemon.Name, pokemon.ID)

	// get pokemon by name - should be cached
	pokemon, _ = client.PokemonByName(context.Background(), "pikachu")

	fmt.Printf("Pokémon: %s (ID: %d)\n", pokemon.Name, pokemon.ID)

	// get pokemon by id
	pokemon, _ = client.PokemonByID(context.Background(), 1)

	fmt.Printf("Pokémon: %s (ID: %d)\n", pokemon.Name, pokemon.ID)

	// get all pokemon
	allPokemon, _ := client.AllPokemon(context.Background())

	fmt.Printf("Total pokemon: %d\n", len(allPokemon))
	fmt.Printf("First Pokémon: %s (ID: %d)\n", allPokemon[0].Name, allPokemon[0].ID)
	fmt.Printf("Last Pokémon: %s (ID: %d)\n",
		allPokemon[len(allPokemon)-1].Name, allPokemon[len(allPokemon)-1].ID)

	// get generation by id
	generation, _ := client.GenerationByID(context.Background(), 1)

	fmt.Printf("Generation: %s (ID: %d)\n", generation.Name, generation.ID)

	// get generation by name
	generation, _ = client.GenerationByName(context.Background(), "generation-iii")

	fmt.Printf("Generation: %s (ID: %d)\n", generation.Name, generation.ID)

	// get all generations
	generations, _ := client.AllGenerations(context.Background())

	fmt.Printf("First Generation: %s (ID: %d)\n", generations[0].Name, generations[0].ID)

	// paginating
	pager := client.PokemonPager(pokeapi.Limit(50))
	firstPage, _ := pager.Next(context.Background())
	res, _ := pokeapi.FetchResultsN[pokeapi.Pokemon](context.Background(), client, firstPage, 5)
	fmt.Printf("First Pokemon: %s\n", res[0].Name)

	// next page
	secondPage, _ := pager.Next(context.Background())
	res, _ = pokeapi.FetchResultsN[pokeapi.Pokemon](context.Background(), client, secondPage, 5)
	fmt.Printf("Second Pokemon: %s\n", res[0].Name)

	// previous page
	previousPage, _ := pager.Previous(context.Background())
	res, _ = pokeapi.FetchResultsN[pokeapi.Pokemon](context.Background(), client, previousPage, 5)
	fmt.Printf("Previous Pokemon: %s\n", res[0].Name)
}
