package enums

import (
	"fmt"
	"io"
	"strconv"
)

type SystemeEchange string

const (
	SystemeEchangeDon   SystemeEchange = "Don"
	SystemeEchangeTroc  SystemeEchange = "Troc"
	SystemeEchangeVente SystemeEchange = "Vente"
)

var AllSystemeEchange = []SystemeEchange{
	SystemeEchangeDon,
	SystemeEchangeTroc,
	SystemeEchangeVente,
}

func (e SystemeEchange) IsValid() bool {
	switch e {
	case SystemeEchangeDon, SystemeEchangeTroc, SystemeEchangeVente:
		return true
	}
	return false
}

func (e SystemeEchange) String() string {
	return string(e)
}

func (e *SystemeEchange) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SystemeEchange(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SystemeEchange", str)
	}
	return nil
}

func (e SystemeEchange) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
