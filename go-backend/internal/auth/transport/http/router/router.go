package router

import (
	"purple/internal/auth"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App, h auth.Handler) {
	g := app.Group("/auth")
	g.Post("/signup", h.Signup)
	g.Post("/login", h.Login)
}
