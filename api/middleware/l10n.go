package middleware

import (
	"github.com/gofiber/fiber/v2"
)

var AcceptLanguage = [27]string{
	"ba_BAS", "br_BRE",
	"fr_FR",
	"gr_ALS", "gr_FLA", "gr_FRA",
	"oc_CRO", "oc_GAS", "oc_LAN", "oc_NOC", "oc_PRO",
	"ro_ANG", "ro_BOU", "ro_CAT", "ro_CEN", "ro_CHA",
	"ro_COR", "ro_FRC", "ro_FRP", "ro_GAL", "ro_LIG", "ro_LOR",
	"ro_MAI", "ro_NOR", "ro_PIC", "ro_POI", "ro_WAL",
}

func L10n() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		lang := c.Context().Request.Header.Peek("Content-Language")
		c.Locals("lang", "fr_FR")
		for _, acceptLang := range AcceptLanguage {
			if string(lang) == acceptLang {
				c.Locals("lang", string(lang))
				break
			}
		}
		return c.Next()
	}
}
