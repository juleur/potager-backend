package mailer

import (
	"math/rand"
	"strconv"
	"time"
)

func GenerateCode() string {
	rand.Seed(time.Now().UTC().UnixNano())
	min := 1000
	max := 9999
	code := rand.Intn((max - min) + min)
	return strconv.FormatInt(int64(code), 10)
}
