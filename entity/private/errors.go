package private

import (
	"npp_backend/l10n/translate"
)

type Error struct {
	Location   string //"db.person.person_psql.go"
	Line       int
	Err        error
	TranslKey  translate.Key // FarmerAlreadyMuted
	ErrorCode  int
	StatusCode int // http error
}
