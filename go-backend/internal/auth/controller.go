package auth

import (
	"context"
	"purple/internal/auth/domain"
)

type Controller interface {
	Signup(ctx context.Context, params domain.SignupReq) (domain.AuthResp, error)
	Login(ctx context.Context, params domain.LoginReq) (domain.AuthResp, error)
}
