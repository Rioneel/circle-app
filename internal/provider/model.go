package provider

type Provider struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Area     string `json:"area"`
}

type SearchResult struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	WeightedScore float64 `json:"weightedScore"`
	Voices        int64   `json:"voices"`
}

type Recommendation struct {
	UserID   string  `json:"userId"`
	UserName string  `json:"userName"`
	Rating   float64 `json:"rating"`
}