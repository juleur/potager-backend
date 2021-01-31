package public

type NearestPotager struct {
	User         User      `json:"user"`
	Commune      string    `json:"commune"`
	Coordonnees  []float64 `json:"coordonnees"`
	Favorite     bool      `json:"favorite"`
	FruitsCount  int       `json:"fruitsCount"`
	LegumesCount int       `json:"legumesCount"`
	GrainesCount int       `json:"grainesCount"`
}
