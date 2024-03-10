package ws

import (
	"purple/internal/api"
	"purple/internal/api/transport/ws/middleware"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func SetupSessionSocket(app *fiber.App, h api.ChatHandler, cfg *websocket.Config) {
	g := app.Group("/ws")
	g.Use(middleware.WsProtocolUpgrade())

	g.Get("/session/:id", websocket.New(h.Chat, *cfg))
}
