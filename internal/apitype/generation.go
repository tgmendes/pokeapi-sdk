package apitype

type Generation struct {
	Abilities      []interface{}      `json:"abilities"`
	ID             int                `json:"id"`
	MainRegion     NamedResource      `json:"main_region"`
	Moves          []NamedResource    `json:"moves"`
	Name           string             `json:"name"`
	Names          []LanguageResource `json:"names"`
	PokemonSpecies []NamedResource    `json:"pokemon_species"`
	Types          []NamedResource    `json:"types"`
	VersionGroups  []NamedResource    `json:"version_groups"`
}

type LanguageResource struct {
	Language NamedResource `json:"language"`
	Name     string        `json:"name"`
}
