package token

import (
	"bytes"
	"testing"
)

func TestJwtAuth(t *testing.T) {
	tokens := [][]byte{
		[]byte("Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1MTYyMzkwMjJ9.tbDepxpstvGdW8TC3G8zg4B6rUYAOvfzdceoH48wgRQ"),
		[]byte("jvdfvfdvdvbgfbfg"),
		[]byte("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1MTYyMzkwMjJ9.tbDepxpstvGdW8TC3G8zg4B6rUYAOvfzdceoH48wgRQ"),
		[]byte("Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9eyJpYXQiOjE1MTYyMzkwMjJ9tbDepxpstvGdW8TC3G8zg4B6rUYAOvfzdceoH48wgRQ"),
		[]byte(""),
		[]byte("Bearer JIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1MTYyMzkwMjJ9.tbDepxp"),
		[]byte("Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL2RvbWFpbi5mciIsImV4cCI6MTYxMDg0NTcwNSwiaWF0IjoxNjEwODAyNTA1LCJ1c2VySWQiOjIsInVzZXJuYW1lIjoibWF4MjM0Iiwic3RhdHVzIjpmYWxzZX0.grJw6EqWeJZb2N_ljsGSFMXWN7ZGXFnNUu4dHM1rXSY"),
	}
	for _, bearerJwt := range tokens {
		if match := IsItAJwtToken(bearerJwt); match {
			jwToken := bytes.Split(bearerJwt, []byte(" "))[1]
			userID, err := VerifyJWT(jwToken)
			if err != nil {
				t.Error(err.Err)
			}
			t.Logf("UserID: %d\n", userID)
		} else {
			t.Log("Cette chaîne de caractères n'est pas un JWT Token")
		}
	}
}
