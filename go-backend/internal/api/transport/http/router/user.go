package router

import (
	"github.com/gofiber/fiber/v2"
	"hack/internal/api"
)

func SetupUserRoutes(app fiber.Router, h api.UserHandler, mw api.Middleware) {
	g := app.Group("/users")
	g.Get("/me", mw.Jwt(), h.Me)
	g.Post("/reset_password_code", h.ResetPasswordCode)
	g.Post("/validate_reset_password_code", h.ValidateResetPasswordCode)
	g.Put("/reset_password", h.ResetPassword)
}
