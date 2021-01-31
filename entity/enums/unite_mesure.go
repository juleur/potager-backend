package enums

import (
	"fmt"
	"io"
	"strconv"
)

type UniteMesure string

const (
	UniteMesureBotte UniteMesure = "Botte"
	UniteMesureKg    UniteMesure = "Kg"
	UniteMesurePiece UniteMesure = "Piece"
)

var AllUniteMesure = []UniteMesure{
	UniteMesureBotte,
	UniteMesureKg,
	UniteMesurePiece,
}

func (e UniteMesure) IsValid() bool {
	switch e {
	case UniteMesureBotte, UniteMesureKg, UniteMesurePiece:
		return true
	}
	return false
}

func (e UniteMesure) String() string {
	return string(e)
}

func (e *UniteMesure) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = UniteMesure(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid UniteMesure", str)
	}
	return nil
}

func (e UniteMesure) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
