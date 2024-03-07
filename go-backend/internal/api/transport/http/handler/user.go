package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"hack/internal/api"
	"hack/internal/api/domain"
	"hack/internal/shared"
	"net/http"
)

type userHandler struct {
	controller api.UserController
	validator  *validator.Validate
}

func NewUserHandler(controller api.UserController) api.UserHandler {
	return &userHandler{
		controller: controller,
		validator:  validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (h *userHandler) Me(ctx *fiber.Ctx) error {
	userId := ctx.Locals(shared.UserIdField).(int64)
	user, err := h.controller.Me(ctx.Context(), userId)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(user)
}

func (h *userHandler) ResetPasswordCode(ctx *fiber.Ctx) error {
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

func (h *userHandler) ValidateResetPasswordCode(ctx *fiber.Ctx) error {
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

func (h *userHandler) ResetPassword(ctx *fiber.Ctx) error {
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
