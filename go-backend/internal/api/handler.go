package api

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type ChatHandler interface {
	Chat(c *websocket.Conn)
}

type SessionHandler interface {
	FindOneById(ctx *fiber.Ctx) error
	List(ctx *fiber.Ctx) error
}
