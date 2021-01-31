package handler

import (
	"net/http"
	"npp_backend/api/middleware"
	"npp_backend/db/person"
	"npp_backend/l10n/translate"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/loctools/go-l10n/loc"
)

func addFavoritePotager(personRepo *person.PersonPsql, l10n *loc.Pool) func(c *fiber.Ctx) error {
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
		if err := personRepo.AddFavoritePotager(c.Context(), userId, farmerId); err != nil {
			return c.Status(err.StatusCode).JSON(&fiber.Map{
				"error_code": 1,
				"message":    l10n.GetContext(lng).Tr(err.TranslKey.String()),
			})
		}
		return c.SendStatus(http.StatusAccepted)
	}
}

func addMutedPotager(personRepo *person.PersonPsql, l10n *loc.Pool) func(c *fiber.Ctx) error {
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
		if err := personRepo.AddMutedPotager(c.Context(), userId, farmerId); err != nil {
			return c.Status(err.StatusCode).JSON(&fiber.Map{
				"error_code": 1,
				"message":    l10n.GetContext(lng).Tr(err.TranslKey.String()),
			})
		}
		return c.SendStatus(http.StatusAccepted)
	}
}

func addNewFarmer(personRepo *person.PersonPsql, l10n *loc.Pool) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func createFruits(personRepo *person.PersonPsql, l10n *loc.Pool) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func createGraines(personRepo *person.PersonPsql, l10n *loc.Pool) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func createLegumes(personRepo *person.PersonPsql, l10n *loc.Pool) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func deleteFruit(personRepo *person.PersonPsql, l10n *loc.Pool) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		lng := c.Locals("lang").(string)

		fruitIdStr := c.Params("fruit_id")
		fruitId, err := strconv.Atoi(fruitIdStr)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
				"error_code": 1,
				"message":    l10n.GetContext(lng).Tr(translate.KeyInternalServerError.String()),
			})
		}
		if err := personRepo.DeleteFruit(c.Context(), fruitId); err != nil {
			return c.Status(err.StatusCode).JSON(&fiber.Map{
				"error_code": 1,
				"message":    l10n.GetContext(lng).Tr(err.TranslKey.String()),
			})
		}
		return c.SendStatus(http.StatusAccepted)
	}
}

func deleteGraine(personRepo *person.PersonPsql, l10n *loc.Pool) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		lng := c.Locals("lang").(string)

		graineIdStr := c.Params("graine_id")
		graineId, err := strconv.Atoi(graineIdStr)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
				"error_code": 1,
				"message":    l10n.GetContext(lng).Tr(translate.KeyInternalServerError.String()),
			})
		}
		if err := personRepo.DeleteGraine(c.Context(), graineId); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
				"error_code": 1,
				"message":    l10n.GetContext(lng).Tr(translate.KeyInternalServerError.String()),
			})
		}
		return c.SendStatus(http.StatusAccepted)
	}
}

func deleteLegume(personRepo *person.PersonPsql, l10n *loc.Pool) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		lng := c.Locals("lang").(string)

		legumeIdStr := c.Params("graine_id")
		legumeId, err := strconv.Atoi(legumeIdStr)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
				"error_code": 1,
				"message":    l10n.GetContext(lng).Tr(translate.KeyInternalServerError.String()),
			})
		}
		if err := personRepo.DeleteLegume(c.Context(), legumeId); err != nil {
			return c.Status(err.StatusCode).JSON(&fiber.Map{
				"error_code": 1,
				"message":    l10n.GetContext(lng).Tr(err.TranslKey.String()),
			})
		}
		return c.SendStatus(http.StatusAccepted)
	}
}

func updateFarmerInfo(personRepo *person.PersonPsql, l10n *loc.Pool) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func updateFruit(personRepo *person.PersonPsql, l10n *loc.Pool) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func updateGraine(personRepo *person.PersonPsql, l10n *loc.Pool) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func updateLegume(personRepo *person.PersonPsql, l10n *loc.Pool) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func removeFavoritePotager(personRepo *person.PersonPsql, l10n *loc.Pool) func(c *fiber.Ctx) error {
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
		if err := personRepo.DeleteFavoritePotager(c.Context(), userId, farmerId); err != nil {
			return c.Status(err.StatusCode).JSON(&fiber.Map{
				"error_code": 1,
				"message":    l10n.GetContext(lng).Tr(err.TranslKey.String()),
			})
		}
		return c.SendStatus(http.StatusAccepted)
	}
}

func removeMutedPotager(personRepo *person.PersonPsql, l10n *loc.Pool) func(c *fiber.Ctx) error {
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
		if err := personRepo.DeleteMutedPotager(c.Context(), userId, farmerId); err != nil {
			return c.Status(err.StatusCode).JSON(&fiber.Map{
				"error_code": 1,
				"message":    l10n.GetContext(lng).Tr(err.TranslKey.String()),
			})
		}
		return c.SendStatus(http.StatusAccepted)
	}
}

func MakePersonHandlers(app *fiber.App, personRepo *person.PersonPsql, l10n *loc.Pool) {
	userGroup := app.Group("/u", middleware.Authorization(l10n))
	userGroup.Post("/add_favorite_potager/:farmer_id", addFavoritePotager(personRepo, l10n))
	userGroup.Post("/add_muted_potager/:farmer_id", addMutedPotager(personRepo, l10n))
	userGroup.Post("/add_new_farmer", addNewFarmer(personRepo, l10n))
	userGroup.Delete("/del_fruit/:fruit_id", deleteFruit(personRepo, l10n))
	userGroup.Delete("/del_graine/:graine_id", deleteGraine(personRepo, l10n))
	userGroup.Delete("/del_legume/:legume_id", deleteLegume(personRepo, l10n))
	userGroup.Post("/new_fruits", createFruits(personRepo, l10n))
	userGroup.Post("/new_graines", createGraines(personRepo, l10n))
	userGroup.Post("/new_legumes", createLegumes(personRepo, l10n))
	userGroup.Put("/update_farmer", updateFarmerInfo(personRepo, l10n))
	userGroup.Put("/update_fruit", updateFruit(personRepo, l10n))
	userGroup.Put("/update_graine", updateGraine(personRepo, l10n))
	userGroup.Put("/update_legume", updateLegume(personRepo, l10n))
	userGroup.Delete("/rm_favorite_potager/:farmer_id", removeFavoritePotager(personRepo, l10n))
	userGroup.Delete("/rm_muted_potager/:farmer_id", removeMutedPotager(personRepo, l10n))
}
