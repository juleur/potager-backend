package public

type MutedPotager struct {
	User         User   `json:"user"`
	Farmer       Farmer `json:"farmer"`
	FruitsCount  int    `json:"fruitsCount"`
	LegumesCount int    `json:"legumesCount"`
	GrainesCount int    `json:"grainesCount"`
}
