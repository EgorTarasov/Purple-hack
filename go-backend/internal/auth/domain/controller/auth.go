package controller

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"github.com/yogenyslav/logger"
	"hack/internal/api"
	"hack/internal/api/data"
	"hack/internal/auth"
	"hack/internal/auth/domain"
	"hack/internal/shared"
	"hack/pkg"
	"hack/pkg/jwt"
	"hack/pkg/secure"
)

type authController struct {
	userRepo  api.UserRepo
	tokenRepo auth.TokenRepo
	cfg       *jwt.Config
}

func New(ur api.UserRepo, tr auth.TokenRepo, cfg *jwt.Config) auth.Controller {
	return &authController{
		userRepo:  ur,
		tokenRepo: tr,
		cfg:       cfg,
	}
}

func (ac *authController) Signup(ctx context.Context, params domain.SignupReq) (domain.AuthResp, error) {
	var resp domain.AuthResp
	hashedPassword, err := secure.GetPasswordHash(params.Password)
	if err != nil {
		return resp, shared.ErrPasswordTooLong
	}

	user := data.User{
		Email:     params.Email,
		Password:  hashedPassword,
		FirstName: params.FirstName,
		LastName:  params.LastName,
		CreatedAt: pkg.GetLocalTime(),
		UpdatedAt: pkg.GetLocalTime(),
	}

	userId, err := ac.userRepo.InsertOne(ctx, user)
	if err != nil {
		if pkg.CheckErrDuplicateKey(err) {
			return resp, shared.ErrDuplicateKey
		}
		return resp, err
	}

	if err = ac.tokenRepo.SetIdRequest(ctx, userId, 1); err != nil {
		return resp, err
	}

	resp, err = jwt.AuthenticateUser(jwt.AuthUserParams{
		UserId:         userId,
		HashedPassword: hashedPassword,
		PlainPassword:  params.Password,
		RequestId:      1,
		Jwt:            ac.cfg,
	})
	return resp, err
}

func (ac *authController) Login(ctx context.Context, params domain.LoginReq) (domain.AuthResp, error) {
	var resp domain.AuthResp
	user, err := ac.userRepo.FindOneByEmail(ctx, params.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return resp, shared.ErrNoSuchUser
		}
		logger.Error(err)
		return resp, shared.ErrFindRecord
	}

	cnt, err := ac.tokenRepo.GetIdRequest(ctx, user.Id)
	if err != nil && !errors.Is(err, redis.Nil) {
		return resp, err
	}

	if err = ac.tokenRepo.SetIdRequest(ctx, user.Id, cnt+1); err != nil {
		return resp, err
	}

	resp, err = jwt.AuthenticateUser(jwt.AuthUserParams{
		UserId:         user.Id,
		HashedPassword: user.Password,
		PlainPassword:  params.Password,
		RequestId:      cnt + 1,
		Jwt:            ac.cfg,
	})
	return resp, err
}
