package public

type NearestAliment struct {
	User    User    `json:"user"`
	Farmer  Farmer  `json:"farmer"`
	Aliment Aliment `json:"aliment"`
}
