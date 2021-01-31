package middleware

import (
	"bytes"
	"net/http"
	"npp_backend/l10n/translate"
	"npp_backend/pkg/token"

	"github.com/gofiber/fiber/v2"
	"github.com/loctools/go-l10n/loc"
)

func Authorization(l10n *loc.Pool) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		lng := c.Locals("lang").(string)
		bearerJwt := c.Context().Request.Header.Peek("Authorization")
		if token.IsItAJwtToken(bearerJwt) {
			jwToken := bytes.Split(bearerJwt, []byte(" "))[1]
			userId, err := token.VerifyJWT(jwToken)
			if err != nil {
				return c.Status(err.StatusCode).JSON(&fiber.Map{
					"error_code": err.ErrorCode,
					"message":    l10n.GetContext(lng).Tr(err.TranslKey.String()),
				})
			}
			c.Locals("userId", userId)
		} else {
			return c.Status(http.StatusUnauthorized).JSON(&fiber.Map{
				"error_code": 40,
				"message":    l10n.GetContext(lng).Tr(translate.KeyInternalServerError.String()),
			})
		}
		return c.Next()
	}
}
