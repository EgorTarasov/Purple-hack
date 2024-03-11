package auth

import (
	"context"

	"purple/internal/api/data"
)

type UserRepo interface {
	FindOneByEmail(ctx context.Context, email string) (data.User, error)
}
