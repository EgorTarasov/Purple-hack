package controller

import (
	"context"
	"errors"

	"purple/internal/auth"
	"purple/internal/auth/domain"
	"purple/internal/shared"

	"github.com/jackc/pgx/v5"
	"github.com/yogenyslav/logger"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	repo auth.UserRepo
}

func NewAuthController(ur auth.UserRepo) *AuthController {
	return &AuthController{
		repo: ur,
	}
}

func (ac *AuthController) Login(ctx context.Context, params domain.LoginReq) (int64, error) {
	userDb, err := ac.repo.FindOneByEmail(ctx, params.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, shared.ErrNoSuchUser
		}
		logger.Error(err)
		return 0, shared.ErrFindRecord
	}

	if err = bcrypt.CompareHashAndPassword([]byte(userDb.Password), []byte(params.Password)); err != nil {
		return 0, shared.ErrInvalidPassword
	}
	return userDb.Id, nil
}
