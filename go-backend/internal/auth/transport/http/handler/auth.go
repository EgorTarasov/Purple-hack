package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"hack/internal/auth"
	"hack/internal/auth/domain"
	"hack/pkg/jwt"
	"net/http"
)

type authHandler struct {
	controller auth.Controller
	validator  *validator.Validate
}

func New(controller auth.Controller) auth.Handler {
	return &authHandler{
		controller: controller,
		validator:  validator.New(validator.WithRequiredStructEnabled()),
	}
}

// Signup godoc
//
//	@Summary		Signup
//	@Description	Adds new user to DB; email and password are required
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			signupReq		body		domain.SignupReq		true	"User's data to signup"
//	@Success		201				{object}	string					"User successfully created"
//	@Failure		400				{object}	response.ErrorResponse	"Email is already used"
//	@Failure		400				{object}	response.ErrorResponse	"Password is too long"
//	@Failure		400				{object}	response.ErrorResponse	"Password is too long"
//	@Failure		422				{object}	response.ErrorResponse	"Body validation error"
//	@Router			/auth/signup	[post]
func (h *authHandler) Signup(ctx *fiber.Ctx) error {
	var req domain.SignupReq
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}
	if err := h.validator.Struct(&req); err != nil {
		return err
	}

	resp, err := h.controller.Signup(ctx.Context(), req)
	if err != nil {
		return err
	}

	cookie := jwt.AccessTokenCookie()
	cookie.Value = resp.AccessToken
	ctx.Cookie(cookie)
	return ctx.SendStatus(http.StatusCreated)
}

// Login godoc
//
//	@Summary		Login
//	@Description	Validate credentials and set cookie if valid
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			loginReq	body		domain.LoginReq			true	"User's credentials"
//	@Success		200			{object}	string					"User successfully logged in"
//	@Failure		401			{object}	response.ErrorResponse	"No such user"
//	@Failure		401			{object}	response.ErrorResponse	"Invalid password"
//	@Failure		422			{object}	response.ErrorResponse	"Body validation error"
//	@Router			/auth/login	[post]
func (h *authHandler) Login(ctx *fiber.Ctx) error {
	var req domain.LoginReq
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}
	if err := h.validator.Struct(&req); err != nil {
		return err
	}

	resp, err := h.controller.Login(ctx.Context(), req)
	if err != nil {
		return err
	}

	cookie := jwt.AccessTokenCookie()
	cookie.Value = resp.AccessToken
	ctx.Cookie(cookie)
	return ctx.SendStatus(http.StatusOK)
}
