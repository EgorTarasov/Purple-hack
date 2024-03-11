package controller

import (
	"context"

	"purple/internal/api"
	"purple/internal/shared"
	"purple/pkg"

	"github.com/google/uuid"
	"github.com/yogenyslav/logger"
)

type UserController struct {
	repo api.UserRepo
}

func NewUserController(ur api.UserRepo) *UserController {
	return &UserController{
		repo: ur,
	}
}

func (uc *UserController) SaveSession(ctx context.Context, userId int64, sessionId uuid.UUID) error {
	if err := uc.repo.SaveSession(ctx, userId, sessionId); err != nil {
		if pkg.CheckErrDuplicateKey(err) {
			return nil
		}
		logger.Error(err)
		return shared.ErrInsertRecord
	}
	return nil
}
