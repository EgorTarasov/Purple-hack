package router

import (
	"purple/internal/api"

	"github.com/gofiber/fiber/v2"
)

func SetupSessionRoutes(app *fiber.App, h api.SessionHandler) {
	sessions := app.Group("/api/sessions")
	sessions.Get("/:id", h.FindOneById)
	sessions.Get("/list", h.List)
}
