package token

import (
	"fmt"
	"math/rand"
	"net/http"
	"npp_backend/entity/private"
	"npp_backend/entity/public"
	"npp_backend/l10n/translate"
	"regexp"
	"strings"
	"time"

	"github.com/gbrlsnchs/jwt/v3"
)

const regex = `^Bearer [A-Za-z0-9-_=]+\.[A-Za-z0-9-_=]+\.?[A-Za-z0-9-_.+/=]*$`

var rxJWT = regexp.MustCompile(regex)

var SECRET_KEY = []byte("aj50c0nKdw3EcVReGMWSrIv6vBktkznm")

type JWTPayload struct {
	jwt.Payload
	UserID   int    `json:"userId"`
	Username string `json:"username"`
	Status   bool   `json:"status"` // si utilisateur est vérifié
}

func IsItAJwtToken(jwt []byte) bool {
	return rxJWT.Match(jwt)
}

func VerifyJWT(jwToken []byte) (int, *private.Error) {
	pl := JWTPayload{}
	// revoir comment passer la clé dans le context
	signature := jwt.NewHS256(SECRET_KEY)
	expValidator := jwt.ExpirationTimeValidator(time.Now())
	validatePayload := jwt.ValidatePayload(&pl.Payload, expValidator)
	if _, err := jwt.Verify(jwToken, signature, &pl, validatePayload); err != nil {
		if err == jwt.ErrExpValidation {
			return 0, &private.Error{
				Location:   "pkg.token.VerifyJWT",
				Line:       41,
				Err:        fmt.Errorf("%w: JWT => %s", err, string(jwToken)),
				TranslKey:  translate.KeyInvalidToken,
				ErrorCode:  42,
				StatusCode: http.StatusUnauthorized,
			}
		}
		return 0, &private.Error{
			Location:   "pkg.token.VerifyJWT",
			Line:       39,
			Err:        fmt.Errorf("%w: JWT => %s", err, string(jwToken)),
			TranslKey:  translate.KeyInvalidToken,
			ErrorCode:  41,
			StatusCode: http.StatusUnauthorized,
		}
	}
	return pl.UserID, nil
}

func GenerateTokens(user *private.User) (*public.Tokens, *private.Error) {
	pl := JWTPayload{
		Payload: jwt.Payload{
			Issuer:         "https://domain.fr",
			ExpirationTime: jwt.NumericDate(time.Now().Add(20 * time.Second)),
			IssuedAt:       jwt.NumericDate(time.Now()),
		},
		Username: user.Username,
		UserID:   user.ID,
		Status:   user.IsVerified(),
	}
	jwtToken, err := jwt.Sign(&pl, jwt.NewHS256(SECRET_KEY))
	if err != nil {
		return nil, &private.Error{
			Location:   "pkg.token.GenerateTokens",
			Line:       64,
			Err:        fmt.Errorf("Impossible d'encoder nouveau jwt: user n°%d", user.ID),
			TranslKey:  translate.KeyInternalServerError,
			ErrorCode:  1,
			StatusCode: http.StatusInternalServerError,
		}
	}
	tokens := public.Tokens{}
	tokens.JWT = string(jwtToken)
	tokens.RefreshToken = HexKeyGenerator(32)
	return &tokens, nil
}

func HexKeyGenerator(nb int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const letterBytes = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	sb := strings.Builder{}
	sb.Grow(nb)
	for ; nb > 0; nb-- {
		sb.WriteByte(letterBytes[rand.Intn(len(letterBytes)-1)])
	}
	return sb.String()
}

func IsAlphanumeric(token []byte) bool {
	for _, c := range token {
		if (c >= 48 && c <= 57) || (c >= 65 && c <= 90) || (c >= 97 && c <= 122) {
			continue
		} else {
			return false
		}
	}
	return true
}
