package auth

import "github.com/gofiber/fiber/v2"

type Handler interface {
	Signup(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
}
