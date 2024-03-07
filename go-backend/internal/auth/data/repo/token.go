package repo

import (
	"context"
	"fmt"
	"purple/internal/auth"
	"purple/internal/shared"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type tokenRepo struct {
	rc *redis.Client
}

func NewTokenRepo(rc *redis.Client) auth.TokenRepo {
	return &tokenRepo{
		rc: rc,
	}
}

func (r *tokenRepo) SetIdRequest(ctx context.Context, userId, cnt int64) error {
	return r.rc.Set(ctx, fmt.Sprintf("token-%d", userId), cnt, shared.JwtTokenExp).Err()
}

func (r *tokenRepo) GetIdRequest(ctx context.Context, userId int64) (int64, error) {
	var (
		res       int64
		cntString string
		err       error
	)

	cntString, err = r.rc.Get(ctx, fmt.Sprintf("token-%d", userId)).Result()
	if err != nil {
		return res, err
	}

	res, err = strconv.ParseInt(cntString, 10, 64)
	return res, err
}
