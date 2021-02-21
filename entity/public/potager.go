package public

type Potager struct {
	User    User     `json:"user"`
	Farmer  Farmer   `json:"farmer"`
	Fruits  []Fruit  `json:"fruits"`
	Legumes []Legume `json:"legumes"`
	Graines []Graine `json:"graines"`
}
