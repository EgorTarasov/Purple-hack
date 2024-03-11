package auth

import (
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	Login(ctx *fiber.Ctx) error
	Logout(ctx *fiber.Ctx) error
}
