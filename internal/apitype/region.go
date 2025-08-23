package apitype

type Region struct {
	ID             int                `json:"id"`
	Locations      []NamedResource    `json:"locations"`
	MainGeneration NamedResource      `json:"main_generation"`
	Name           string             `json:"name"`
	Names          []LanguageResource `json:"names"`
	Pokedexes      []NamedResource    `json:"pokedexes"`
	VersionGroups  []NamedResource    `json:"version_groups"`
}
