package public

type Farmer struct {
	ID          int       `json:"id,omitempty"`
	ImgUrl      string    `json:"imgUrl,omitempty"`
	Description string    `json:"description,omitempty"`
	Commune     string    `json:"commune,omitempty"`
	Coordonnees []float64 `json:"coordonnees,omitempty"`
	Favorite    bool      `json:"favorite"`
}
