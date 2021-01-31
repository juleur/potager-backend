package password

import (
	"fmt"
	"testing"
)

func TestHashedPassword(t *testing.T) {
	pwd := "julien"
	hash, err := HashPassword(pwd)
	if err != nil {
		t.Errorf(err.Err.Error())
	}
	fmt.Println(hash)
}
