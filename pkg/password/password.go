package password

import (
	"errors"
	"net/http"
	"npp_backend/entity/private"
	"npp_backend/l10n/translate"

	"github.com/alexedwards/argon2id"
)

func ComparePasswords(password string, hashedPassword string) (bool, *private.Error) {
	if password == "" {
		return false, &private.Error{
			Location:   "pkg.password.ComparePasswords",
			Line:       13,
			Err:        errors.New("Mot de passe non fourni"),
			TranslKey:  translate.KeyInternalServerError,
			ErrorCode:  1,
			StatusCode: http.StatusInternalServerError,
		}
	}
	if hashedPassword == "" {
		return false, &private.Error{
			Location:   "pkg.password.ComparePasswords",
			Line:       23,
			Err:        errors.New("mot de passe hashé non fourni"),
			TranslKey:  translate.KeyInternalServerError,
			ErrorCode:  1,
			StatusCode: http.StatusInternalServerError,
		}
	}
	match, err := argon2id.ComparePasswordAndHash(password, hashedPassword)
	if err != nil {
		return false, &private.Error{
			Location:   "pkg.password.ComparePasswords",
			Line:       34,
			Err:        errors.New("Erreur lors de la comparaison du mdp et du hashé mdp"),
			TranslKey:  translate.KeyInternalServerError,
			ErrorCode:  1,
			StatusCode: http.StatusInternalServerError,
		}
	}
	return match, nil
}

func HashPassword(password string) (string, *private.Error) {
	hashedPassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", &private.Error{
			Location:   "pkg.password.HashPassword",
			Line:       49,
			Err:        errors.New("Erreur lors du hashing du mdp"),
			TranslKey:  translate.KeyInternalServerError,
			ErrorCode:  1,
			StatusCode: http.StatusInternalServerError,
		}
	}
	return hashedPassword, nil
}
