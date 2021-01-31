package handler

import (
	"fmt"
	"net/http"
	"npp_backend/api/middleware"
	"npp_backend/db/auth"
	"npp_backend/l10n/translate"
	"npp_backend/pkg/password"
	"npp_backend/pkg/token"

	"github.com/gofiber/fiber/v2"
	"github.com/loctools/go-l10n/loc"
)

type inputLogin struct {
	Identifiant string `json:"identifiant"`
	Password    string `json:"password"`
}

type inputNewUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func login(db *auth.AuthPsql, l10n *loc.Pool) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		lng := c.Locals("lang").(string)

		body := new(inputLogin)
		if err := c.BodyParser(body); err != nil {
			return c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
				"error_code": 1,
				"message":    l10n.GetContext(lng).Tr(translate.KeyInternalServerError.String()),
			})
		}

		user, err := db.GetUser(c.Context(), body.Identifiant)
		if err != nil {
			return c.Status(err.StatusCode).JSON(&fiber.Map{
				"error_code": err.ErrorCode,
				"message":    l10n.GetContext(lng).Tr(err.TranslKey.String()),
			})
		}
		// check account state
		if user.State != "OK" {
			if user.State == "LOCKED" {
				return c.Status(http.StatusForbidden).JSON(&fiber.Map{
					"error_code": 1,
					"message":    l10n.GetContext(lng).Tr(translate.KeyUserSuspended.String()) + " " + user.LockedUntil(),
				})
			}
			return c.Status(http.StatusUnavailableForLegalReasons).JSON(&fiber.Map{
				"error_code": 1,
				"message":    l10n.GetContext(lng).Tr(translate.KeyUserBanned.String()),
			})
		}

		matched, err := password.ComparePasswords(body.Password, user.Password)
		if !matched {
			if err = db.FailedAttempt(c.Context(), user.ID); err != nil {
				return c.Status(err.StatusCode).JSON(&fiber.Map{
					"error_code": err.ErrorCode,
					"message":    l10n.GetContext(lng).Tr(err.TranslKey.String()),
				})
			}
			return c.Status(http.StatusUnauthorized).JSON(&fiber.Map{
				"error_code": 1,
				"message":    l10n.GetContext(lng).Tr(translate.KeyPasswordNotMatched.String()),
			})
		} else if err != nil {
			return c.Status(err.StatusCode).JSON(&fiber.Map{
				"error_code": err.ErrorCode,
				"message":    l10n.GetContext(lng).Tr(err.TranslKey.String()),
			})
		}

		tokens, err := token.GenerateTokens(user)
		if err != nil {
			return c.Status(err.StatusCode).JSON(&fiber.Map{
				"error_code": err.StatusCode,
				"message":    l10n.GetContext(lng).Tr(err.TranslKey.String()),
			})
		}

		if err = db.UpdateLastLoggedInAt(c.Context(), user.ID); err != nil {
			return c.Status(err.StatusCode).JSON(&fiber.Map{
				"error_code": err.StatusCode,
				"message":    l10n.GetContext(lng).Tr(err.TranslKey.String()),
			})
		}

		if err = db.UpdateRefreshToken(c.Context(), tokens.RefreshToken, user.ID); err != nil {
			return c.Status(err.StatusCode).JSON(&fiber.Map{
				"error_code": err.StatusCode,
				"message":    l10n.GetContext(lng).Tr(err.TranslKey.String()),
			})
		}

		return c.JSON(*tokens)
	}
}

func newUser(db *auth.AuthPsql, l10n *loc.Pool) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		lng := c.Locals("lang").(string)

		body := new(inputNewUser)
		if err := c.BodyParser(body); err != nil {
			return c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
				"error_code": 1,
				"message":    l10n.GetContext(lng).Tr(translate.KeyInternalServerError.String()),
			})
		}

		user, err := db.CreateNewUser(c.Context(), body.Username, body.Email, body.Password)
		if err != nil {
			return c.Status(err.StatusCode).JSON(&fiber.Map{
				"error_code": err.ErrorCode,
				"message":    l10n.GetContext(lng).Tr(err.TranslKey.String()),
			})
		}

		tokens, err := token.GenerateTokens(user)
		if err != nil {
			return c.Status(err.StatusCode).JSON(&fiber.Map{
				"error_code": err.ErrorCode,
				"message":    l10n.GetContext(lng).Tr(err.TranslKey.String()),
			})
		}

		return c.JSON(*tokens)
	}
}

func refreshJWT(db *auth.AuthPsql, l10n *loc.Pool) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		lng := c.Locals("lang").(string)

		refToken := c.Context().Request.Header.Peek("X-Refresh-Token")

		fmt.Printf("refToken: %s\n", string(refToken))

		if len(refToken) != 32 || !token.IsAlphanumeric(refToken) {
			return c.Status(http.StatusUnauthorized).JSON(&fiber.Map{
				"error_code": 1,
				"message":    l10n.GetContext(lng).Tr(translate.KeyInvalidToken.String()),
			})
		}

		user, err := db.IsRefreshTokenExists(c.Context(), string(refToken))
		if err != nil {
			fmt.Println("error")
			fmt.Println(err)
			return c.Status(err.StatusCode).JSON(&fiber.Map{
				"error_code": err.ErrorCode,
				"message":    l10n.GetContext(lng).Tr(err.TranslKey.String()),
			})
		} else if user.UserPermission.LockedUntil.Valid {
			return c.Status(http.StatusUnauthorized).JSON(&fiber.Map{
				"error_code": 1,
				"message":    l10n.GetContext(lng).Tr(translate.KeyUserSuspended.String() + " " + user.LockedUntil()),
			})
		}

		newTokens, err := token.GenerateTokens(user)
		if err != nil {
			return c.Status(err.StatusCode).JSON(&fiber.Map{
				"error_code": err.ErrorCode,
				"message":    l10n.GetContext(lng).Tr(err.TranslKey.String()),
			})
		}

		if err = db.UpdateRefreshToken(c.Context(), newTokens.RefreshToken, user.ID); err != nil {
			return c.Status(err.StatusCode).JSON(&fiber.Map{
				"error_code": err.ErrorCode,
				"message":    l10n.GetContext(lng).Tr(err.TranslKey.String()),
			})
		}

		return c.JSON(*newTokens)
	}
}

type Body struct {
	Email string `json:"email"`
}

func resetPassword(db *auth.AuthPsql, l10n *loc.Pool) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// lng := c.Locals("lang").(string)
		// email := c.Params("email")

		return nil
	}
}

func sendResetCodePassword(db *auth.AuthPsql, l10n *loc.Pool) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// lng := c.Locals("lang").(string)
		return nil
	}
}

func sendCodeVerif(db *auth.AuthPsql, l10n *loc.Pool) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// lng := c.Locals("lang").(string)
		return nil
	}
}

func verifyAccount(db *auth.AuthPsql, l10n *loc.Pool) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// lng := c.Locals("lang").(string)
		return nil
	}
}

func MakeAuthHandlers(app *fiber.App, authRepo *auth.AuthPsql, l10n *loc.Pool) {
	authGroup := app.Group("/a")
	authGroup.Post("/login", login(authRepo, l10n))
	authGroup.Post("/new_user", newUser(authRepo, l10n))
	authGroup.Get("/refresh_jwt", refreshJWT(authRepo, l10n), middleware.RefreshToken(l10n))
	authGroup.Post("/reset_pwd/send", resetPassword(authRepo, l10n))
	authGroup.Post("/reset_pwd/confirm", sendResetCodePassword(authRepo, l10n))
	authGroup.Get("/verify_account/send", sendCodeVerif(authRepo, l10n), middleware.Authorization(l10n))            // ajout middleware
	authGroup.Post("/verify_account/confirm/:token", verifyAccount(authRepo, l10n), middleware.Authorization(l10n)) // ajout middleware
}
