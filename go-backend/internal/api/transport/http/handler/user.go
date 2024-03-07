package handler

import (
	"net/http"
	"purple/internal/api"
	"purple/internal/api/domain"
	"purple/internal/shared"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type userHttpHandler struct {
	controller api.UserController
	validator  *validator.Validate
}

func NewUserHttpHandler(controller api.UserController) api.UserHttpHandler {
	return &userHttpHandler{
		controller: controller,
		validator:  validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (h *userHttpHandler) Me(ctx *fiber.Ctx) error {
	userId := ctx.Locals(shared.UserIdField).(int64)
	user, err := h.controller.Me(ctx.Context(), userId)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(user)
}

func (h *userHttpHandler) ResetPasswordCode(ctx *fiber.Ctx) error {
	var (
		req domain.ResetPasswordCodeReq
		err error
	)

	if err = ctx.BodyParser(&req); err != nil {
		return err
	}
	if err = h.validator.Struct(&req); err != nil {
		return err
	}
	if err = h.controller.ResetPasswordCode(ctx.Context(), req.Email); err != nil {
		return err
	}

	return ctx.SendStatus(http.StatusNoContent)
}

func (h *userHttpHandler) ValidateResetPasswordCode(ctx *fiber.Ctx) error {
	var (
		req domain.ValidateResetPasswordReq
		err error
	)

	if err = ctx.BodyParser(&req); err != nil {
		return err
	}
	if err = h.validator.Struct(&req); err != nil {
		return err
	}
	if err = h.controller.ValidateResetPasswordCode(ctx.Context(), req.Email, req.Code); err != nil {
		return err
	}

	return ctx.SendStatus(http.StatusNoContent)
}

func (h *userHttpHandler) ResetPassword(ctx *fiber.Ctx) error {
	var (
		req domain.ResetPasswordReq
		err error
	)

	if err = ctx.BodyParser(&req); err != nil {
		return err
	}
	if err = h.validator.Struct(&req); err != nil {
		return err
	}
	if err = h.controller.ResetPassword(ctx.Context(), req.Email, req.NewPassword); err != nil {
		return err
	}

	return ctx.SendStatus(http.StatusNoContent)
}
