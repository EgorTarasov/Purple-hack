package api

import (
	"context"
	"purple/internal/api/data"
)

type UserRepo interface {
	InsertOne(ctx context.Context, user data.User) (int64, error)
	FindOneByEmail(ctx context.Context, email string) (data.User, error)
	FindOneById(ctx context.Context, userId int64) (data.User, error)
	UpdateLastPassReset(ctx context.Context, userId int64) error
	UpdatePassword(ctx context.Context, userId int64, password string) error

	SetPasswordResetTOTP(ctx context.Context, userId int64, secret string) error
	GetPasswordResetTOTP(ctx context.Context, userId int64) (string, error)
	DeletePasswordResetTOTP(ctx context.Context, userId int64) error
	SetResetPassAvailable(ctx context.Context, userId int64) error
	GetDelResetPassAvailable(ctx context.Context, userId int64) error
}
