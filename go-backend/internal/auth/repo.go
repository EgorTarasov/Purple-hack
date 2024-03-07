package auth

import (
	"context"
)

type TokenRepo interface {
	SetIdRequest(ctx context.Context, userId, cnt int64) error
	GetIdRequest(ctx context.Context, userId int64) (int64, error)
}
