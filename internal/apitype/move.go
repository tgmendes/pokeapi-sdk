package apitype

type Move struct {
	Accuracy           int                `json:"accuracy"`
	ContestCombos      ContestCombos      `json:"contest_combos"`
	ContestEffect      URLResource        `json:"contest_effect"`
	NamedResource      NamedResource      `json:"contest_type"`
	DamageClass        NamedResource      `json:"damage_class"`
	EffectChance       interface{}        `json:"effect_chance"`
	EffectChanges      []interface{}      `json:"effect_changes"`
	EffectEntries      []EffectEntry      `json:"effect_entries"`
	FlavorTextEntries  []FlavorTextEntry  `json:"flavor_text_entries"`
	Generation         NamedResource      `json:"generation"`
	ID                 int                `json:"id"`
	LearnedByPokemon   []NamedResource    `json:"learned_by_pokemon"`
	Machines           []interface{}      `json:"machines"`
	Meta               Meta               `json:"meta"`
	Name               string             `json:"name"`
	Names              []LanguageResource `json:"names"`
	PastValues         []interface{}      `json:"past_values"`
	Power              int                `json:"power"`
	Pp                 int                `json:"pp"`
	Priority           int                `json:"priority"`
	StatChanges        []interface{}      `json:"stat_changes"`
	SuperContestEffect URLResource        `json:"super_contest_effect"`
	Target             NamedResource      `json:"target"`
	Type               NamedResource      `json:"type"`
}

type ContestCombos struct {
	Normal Normal `json:"normal"`
	Super  Normal `json:"super"`
}

type Normal struct {
	UseAfter  interface{}     `json:"use_after"`
	UseBefore []NamedResource `json:"use_before"`
}

type EffectEntry struct {
	Effect      string        `json:"effect"`
	Language    NamedResource `json:"language"`
	ShortEffect string        `json:"short_effect"`
}

type FlavorTextEntry struct {
	FlavorText   string        `json:"flavor_text"`
	Language     NamedResource `json:"language"`
	VersionGroup NamedResource `json:"version_group"`
}

type Meta struct {
	Ailment       NamedResource `json:"ailment"`
	AilmentChance int           `json:"ailment_chance"`
	Category      NamedResource `json:"category"`
	CritRate      int           `json:"crit_rate"`
	Drain         int           `json:"drain"`
	FlinchChance  int           `json:"flinch_chance"`
	Healing       int           `json:"healing"`
	MaxHits       interface{}   `json:"max_hits"`
	MaxTurns      interface{}   `json:"max_turns"`
	MinHits       interface{}   `json:"min_hits"`
	MinTurns      interface{}   `json:"min_turns"`
	StatChance    int           `json:"stat_chance"`
}
