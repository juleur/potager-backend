package handler

import (
	"net/http"
	"npp_backend/api/middleware"
	"npp_backend/db/farmer"

	"npp_backend/l10n/translate"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/loctools/go-l10n/loc"
)

func fetchFavoritePotagers(farmerRepo *farmer.FarmerPsql, l10n *loc.Pool) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		lng := c.Locals("lang").(string)

		userIdStr := c.Params("user_id")
		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
				"error_code": 1,
				"message":    l10n.GetContext(lng).Tr(translate.KeyInvalidToken.String()),
			})
		}

		favorites, er := farmerRepo.GetFavoritePotagers(c.Context(), userId)
		if er != nil {
			return c.Status(er.StatusCode).JSON(&fiber.Map{
				"error_code": 1,
				"message":    l10n.GetContext(lng).Tr(er.TranslKey.String()),
			})
		}

		return c.JSON(favorites)
	}
}

func fetchMutedPotagers(farmerRepo *farmer.FarmerPsql, l10n *loc.Pool) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		lng := c.Locals("lang").(string)

		userIdStr := c.Params("user_id")
		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
				"error_code": 1,
				"message":    l10n.GetContext(lng).Tr(translate.KeyInternalServerError.String()),
			})
		}

		mutes, er := farmerRepo.GetMutedPotagers(c.Context(), userId)
		if er != nil {
			return c.Status(er.StatusCode).JSON(&fiber.Map{
				"error_code": er.ErrorCode,
				"message":    l10n.GetContext(lng).Tr(er.TranslKey.String()),
			})
		}

		return c.JSON(mutes)
	}
}

func fetchPotager(farmerRepo *farmer.FarmerPsql, l10n *loc.Pool) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		lng := c.Locals("lang").(string)
		userId := c.Locals("userId").(int)

		farmerIdStr := c.Params("farmer_id")
		farmerId, err := strconv.Atoi(farmerIdStr)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
				"error_code": 1,
				"message":    l10n.GetContext(lng).Tr(translate.KeyInternalServerError.String()),
			})
		}

		potager, er := farmerRepo.GetPotager(c.Context(), userId, farmerId)
		if er != nil {
			return c.Status(er.StatusCode).JSON(&fiber.Map{
				"error_code": er.ErrorCode,
				"message":    l10n.GetContext(lng).Tr(er.TranslKey.String()),
			})
		}

		return c.JSON(*potager)
	}
}

type GeoSearch struct {
	Coordonnees []float64 `json:"coordonnees"`
	Search      string    `json:"search"`
}

func findNearestAliments(farmerRepo *farmer.FarmerPsql, l10n *loc.Pool) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		lng := c.Locals("lang").(string)
		userId := c.Locals("userId").(int)

		geoBody := new(GeoSearch)
		if err := c.BodyParser(geoBody); err != nil {
			return c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
				"error_code": 1,
				"message":    l10n.GetContext(lng).Tr(translate.KeyInternalServerError.String()),
			})
		}

		nearestAliments, err := farmerRepo.GetNearestAliments(c.Context(), userId, geoBody.Coordonnees, geoBody.Search)
		if err != nil {
			return c.Status(err.StatusCode).JSON(&fiber.Map{
				"error_code": 1,
				"message":    l10n.GetContext(lng).Tr(err.TranslKey.String()),
			})
		}

		return c.JSON(nearestAliments)
	}
}

type Geolocation struct {
	Coordonnees []float64 `json:"coordonnees"`
}

func findNearestPotagers(farmerRepo *farmer.FarmerPsql, l10n *loc.Pool) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		lng := c.Locals("lang").(string)
		userId := c.Locals("userId").(int)

		geoBody := new(Geolocation)
		if err := c.BodyParser(geoBody); err != nil {
			return c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
				"error_code": 1,
				"message":    l10n.GetContext(lng).Tr(translate.KeyInternalServerError.String()),
			})
		}

		nearestPotagers, err := farmerRepo.GetNearestPotagers(c.Context(), userId, geoBody.Coordonnees)
		if err != nil {
			return c.Status(err.StatusCode).JSON(&fiber.Map{
				"error_code": 1,
				"message":    l10n.GetContext(lng).Tr(err.TranslKey.String()),
			})
		}

		return c.JSON(nearestPotagers)
	}
}

func MakeFarmerHandlers(app *fiber.App, farmerRepo *farmer.FarmerPsql, l10n *loc.Pool) {
	farmerGroup := app.Group("/f", middleware.Authorization(l10n))
	farmerGroup.Get("/fetch_favorite/:user_id", fetchFavoritePotagers(farmerRepo, l10n))
	farmerGroup.Get("/fetch_muted/:user_id", fetchMutedPotagers(farmerRepo, l10n))
	farmerGroup.Get("/fetch_potager/:farmer_id", fetchPotager(farmerRepo, l10n))
	farmerGroup.Post("/find_nearest_aliments", findNearestAliments(farmerRepo, l10n))
	farmerGroup.Post("/find_nearest_potagers", findNearestPotagers(farmerRepo, l10n))
}
