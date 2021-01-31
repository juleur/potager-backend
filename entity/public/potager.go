package public

type Potager struct {
	User    User     `json:"user"`
	Fruits  []Fruit  `json:"fruits"`
	Legumes []Legume `json:"legumes"`
	Graines []Graine `json:"graines"`
}
