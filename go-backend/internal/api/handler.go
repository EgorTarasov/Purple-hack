package api

import (
	"github.com/gofiber/contrib/websocket"
)

type SessionWsHandler interface {
	Serve()
	Chat(c *websocket.Conn)
}
