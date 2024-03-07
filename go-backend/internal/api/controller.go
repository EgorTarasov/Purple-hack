package api

import (
	"context"
	"hack/internal/api/domain"
)

type UserController interface {
	Me(ctx context.Context, userId int64) (domain.User, error)
	ResetPasswordCode(ctx context.Context, email string) error
	ValidateResetPasswordCode(ctx context.Context, email, code string) error
	ResetPassword(ctx context.Context, email, password string) error
}
