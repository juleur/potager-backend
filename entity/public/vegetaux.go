package public

import (
	"npp_backend/entity/enums"
)

type Fruit struct {
	ID             int                    `json:"id,omitempty"`
	ImgUrl         string                 `json:"imgUrl,omitempty"`
	Nom            string                 `json:"nom,omitempty"`
	Variete        string                 `json:"variete,omitempty"`
	SystemeEchange []enums.SystemeEchange `json:"systemeEchange,omitempty"`
	Prix           float64                `json:"prix,omitempty"`
	UniteMesure    enums.UniteMesure      `json:"uniteMesure,omitempty"`
	Stock          int                    `json:"stock,omitempty"`
}

type Graine struct {
	ID             int                    `json:"id,omitempty"`
	ImgUrl         string                 `json:"imgUrl,omitempty"`
	Nom            string                 `json:"nom,omitempty"`
	Variete        string                 `json:"variete,omitempty"`
	SystemeEchange []enums.SystemeEchange `json:"systemeEchange,omitempty"`
	Prix           float64                `json:"prix,omitempty"`
	Stock          int                    `json:"stock,omitempty"`
}

type Legume struct {
	ID             int                    `json:"id,omitempty"`
	ImgUrl         string                 `json:"imgUrl,omitempty"`
	Nom            string                 `json:"nom,omitempty"`
	Variete        string                 `json:"variete,omitempty"`
	SystemeEchange []enums.SystemeEchange `json:"systemeEchange,omitempty"`
	Prix           float64                `json:"prix,omitempty"`
	UniteMesure    enums.UniteMesure      `json:"uniteMesure,omitempty"`
	Stock          int                    `json:"stock,omitempty"`
}
