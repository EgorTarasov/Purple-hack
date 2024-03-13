package handler

import (
	"net/http"
	"strconv"

	"purple/internal/api"
	"purple/internal/shared"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type SessionHandler struct {
	controller api.SessionController
}

func NewSessionHandler(sc api.SessionController) *SessionHandler {
	return &SessionHandler{
		controller: sc,
	}
}

func (h *SessionHandler) FindOneById(ctx *fiber.Ctx) error {
	sessionId, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return shared.ErrSessionIdInvalid
	}

	session, err := h.controller.FindOneById(ctx.Context(), sessionId)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(session)
}

func (h *SessionHandler) List(ctx *fiber.Ctx) error {
	userId, err := strconv.ParseInt(ctx.Cookies("auth"), 10, 64)
	if err != nil {
		return shared.ErrUserIdInvalid
	}

	sessions, err := h.controller.List(ctx.Context(), userId)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(sessions)
}
