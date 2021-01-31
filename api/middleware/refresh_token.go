package middleware

import (
	"net/http"
	"npp_backend/l10n/translate"
	"npp_backend/pkg/token"

	"github.com/gofiber/fiber/v2"
	"github.com/loctools/go-l10n/loc"
)

func RefreshToken(l10n *loc.Pool) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		lng := c.Locals("lang").(string)

		refToken := c.Context().Request.Header.Peek("X-Refresh-Token")
		if len(refToken) != 32 || !token.IsAlphanumeric(refToken) {
			return c.Status(http.StatusUnauthorized).JSON(&fiber.Map{
				"error_code": 41,
				"message":    l10n.GetContext(lng).Tr(translate.KeyInvalidToken.String()),
			})
		}
		return c.Next()
	}
}
