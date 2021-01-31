package public

import (
	"npp_backend/entity/enums"
)

type Aliment struct {
	ID             int                    `json:"id"`
	ImgUrl         string                 `json:"imgUrl"`
	Nom            string                 `json:"nom"`
	Variete        string                 `json:"variete"`
	SystemeEchange []enums.SystemeEchange `json:"systemeEchange"`
	Prix           float64                `json:"prix"`
	UniteMesure    enums.UniteMesure      `json:"uniteMesure"`
	Stock          int                    `json:"stock"`
}
