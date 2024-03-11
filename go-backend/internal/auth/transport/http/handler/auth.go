package handler

import (
	"net/http"
	"strconv"
	"time"

	"purple/internal/auth"
	"purple/internal/auth/domain"
	"purple/pkg"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	controller auth.Controller
	validator  *validator.Validate
}

func NewAuthHandler(controller auth.Controller) *AuthHandler {
	return &AuthHandler{
		controller: controller,
		validator:  validator.New(validator.WithRequiredStructEnabled()),
	}
}

// Login godoc
//
//	@Summary		Login
//	@Description	Validate credentials and set cookie if valid
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			loginReq	body	domain.LoginReq	true	"User's credentials"
//	@Success		200			"User successfully logged in"
//	@Failure		401			{object}	response.ErrorResponse	"No such user"
//	@Failure		401			{object}	response.ErrorResponse	"Invalid password"
//	@Failure		422			{object}	response.ErrorResponse	"Body validation error"
//	@Router			/auth/login	[post]
func (h *AuthHandler) Login(ctx *fiber.Ctx) error {
	var req domain.LoginReq
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}
	if err := h.validator.Struct(&req); err != nil {
		return err
	}

	userId, err := h.controller.Login(ctx.Context(), req)
	if err != nil {
		return err
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "auth",
		Value:    strconv.FormatInt(userId, 10),
		Path:     "/",
		MaxAge:   int((time.Hour * 12).Seconds()),
		Expires:  pkg.GetLocalTime().Add(time.Hour * 12),
		Secure:   true,
		HTTPOnly: false,
	})

	return ctx.SendStatus(http.StatusOK)
}

// Logout godoc
//
//	@Summary		Logout
//	@Description	Resets auth cookie
//	@Tags			auth
//	@Success		204				"Successfully logged out"
//	@Router			/auth/logout	[delete]
func (h *AuthHandler) Logout(ctx *fiber.Ctx) error {
	ctx.Cookie(&fiber.Cookie{
		Name:  "auth",
		Value: "",
	})
	return ctx.SendStatus(http.StatusNoContent)
}
