package auth

import (
	"context"

	"purple/internal/auth/domain"
)

type Controller interface {
	Login(ctx context.Context, params domain.LoginReq) (int64, error)
}
