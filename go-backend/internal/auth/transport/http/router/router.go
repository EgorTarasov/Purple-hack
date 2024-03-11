package router

import (
	"purple/internal/auth"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app *fiber.App, h auth.Handler) {
	g := app.Group("/auth")
	g.Post("/login", h.Login)
	g.Delete("/logout", h.Logout)
}
