package api

import (
	"github.com/gofiber/contrib/websocket"
)

type ChatHandler interface {
	Chat(c *websocket.Conn)
}
