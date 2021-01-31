package public

type User struct {
	ID          int       `json:"id"`
	Username    string    `json:"username"`
	ImgURL      string    `json:"img_url,omitempty"`
	Description string    `json:"description,omitempty"`
	Commune     string    `json:"commune,omitempty"`
	Coordonnees []float64 `json:"coordonnees,omitempty"`
}
