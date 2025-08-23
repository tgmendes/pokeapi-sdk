package apitype

type List struct {
	Count   int             `json:"count"`
	Next    *string         `json:"next"`
	Prev    *string         `json:"previous"`
	Results []NamedResource `json:"results"`
}

type NamedResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
