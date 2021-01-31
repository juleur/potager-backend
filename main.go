package main

import (
	"log"
	"npp_backend/api/handler"
	"npp_backend/api/middleware"
	"npp_backend/db"
	"npp_backend/db/auth"
	"npp_backend/db/farmer"
	"npp_backend/db/person"
	"npp_backend/l10n"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// TODO: revoir façon de gérer la clé secret pour génération des jwt
	db := db.OpenDB()
	defer db.Close()

	l10n := l10n.LoadTranslations()

	personRepo := person.NewPersonRepository(db)
	farmerRepo := farmer.NewFarmerRepository(db)
	authRepo := auth.NewAuthRepository(db)

	app := fiber.New()

	app.Use(logger.New()) // à désactiver en prod

	app.Use(middleware.L10n())

	// TODO: passer une struct avec app, l10n, mailer, et hébergeur d'images
	// TODO: Manager{App, L10n, Mailer, ImageUploader}
	handler.MakeFarmerHandlers(app, farmerRepo, l10n)
	handler.MakePersonHandlers(app, personRepo, l10n)
	handler.MakeAuthHandlers(app, authRepo, l10n)

	log.Fatal(app.Listen(":9377"))
}
