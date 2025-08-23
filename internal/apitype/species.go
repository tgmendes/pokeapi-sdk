package apitype

type Species struct {
	BaseHappiness        int                `json:"base_happiness"`
	CaptureRate          int                `json:"capture_rate"`
	Color                NamedResource      `json:"color"`
	EggGroups            []NamedResource    `json:"egg_groups"`
	EvolutionChain       URLResource        `json:"evolution_chain"`
	EvolvesFromSpecies   NamedResource      `json:"evolves_from_species"`
	FlavorTextEntries    []FlavorTextEntry  `json:"flavor_text_entries"`
	FormDescriptions     []interface{}      `json:"form_descriptions"`
	FormsSwitchable      bool               `json:"forms_switchable"`
	GenderRate           int                `json:"gender_rate"`
	Genera               []Genus            `json:"genera"`
	Generation           NamedResource      `json:"generation"`
	GrowthRate           NamedResource      `json:"growth_rate"`
	Habitat              NamedResource      `json:"habitat"`
	HasGenderDifferences bool               `json:"has_gender_differences"`
	HatchCounter         int                `json:"hatch_counter"`
	ID                   int                `json:"id"`
	IsBaby               bool               `json:"is_baby"`
	IsLegendary          bool               `json:"is_legendary"`
	IsMythical           bool               `json:"is_mythical"`
	Name                 string             `json:"name"`
	Names                []LanguageResource `json:"names"`
	Order                int                `json:"order"`
	PalParkEncounters    []PalParkEncounter `json:"pal_park_encounters"`
	PokedexNumbers       []PokedexNumber    `json:"pokedex_numbers"`
	Shape                NamedResource      `json:"shape"`
	Varieties            []Variety          `json:"varieties"`
}

type URLResource struct {
	URL string `json:"url"`
}

type Genus struct {
	Genus    string        `json:"genus"`
	Language NamedResource `json:"language"`
}

type PalParkEncounter struct {
	Area      NamedResource `json:"area"`
	BaseScore int           `json:"base_score"`
	Rate      int           `json:"rate"`
}

type PokedexNumber struct {
	EntryNumber int           `json:"entry_number"`
	Pokedex     NamedResource `json:"pokedex"`
}

type Variety struct {
	IsDefault bool          `json:"is_default"`
	Pokemon   NamedResource `json:"pokemon"`
}
