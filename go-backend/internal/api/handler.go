package api

import "github.com/gofiber/fiber/v2"

type UserHttpHandler interface {
	Me(ctx *fiber.Ctx) error
	ResetPasswordCode(ctx *fiber.Ctx) error
	ValidateResetPasswordCode(ctx *fiber.Ctx) error
	ResetPassword(ctx *fiber.Ctx) error
}
