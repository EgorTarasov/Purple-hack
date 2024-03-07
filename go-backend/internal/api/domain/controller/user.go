package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"github.com/yogenyslav/logger"
	"hack/internal/api"
	"hack/internal/api/domain"
	"hack/internal/shared"
	"hack/pkg"
	"hack/pkg/mailing"
	"hack/pkg/secure"
)

type userController struct {
	userRepo   api.UserRepo
	mailServer *mailing.MailServer
}

func NewUserController(ur api.UserRepo, mailServer *mailing.MailServer) api.UserController {
	return &userController{
		userRepo:   ur,
		mailServer: mailServer,
	}
}

func (uc *userController) Me(ctx context.Context, userId int64) (domain.User, error) {
	var resp domain.User
	u, err := uc.userRepo.FindOneById(ctx, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return resp, shared.ErrNoSuchUser
		}
		logger.Error(err)
		return resp, shared.ErrFindRecord
	}

	return domain.User{
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}, nil
}

func (uc *userController) ResetPasswordCode(ctx context.Context, email string) error {
	var err error

	user, err := uc.userRepo.FindOneByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return shared.ErrNoSuchUser
		}
		return shared.ErrFindRecord
	}

	if pkg.GetLocalTime().Unix()-user.LastPassReset.Unix() < int64(shared.PasswordResetTimeout.Seconds()) {
		return shared.ErrPassRecentlyReset
	}

	code, err := secure.GeneratePassResetCode()
	if err != nil {
		return err
	}
	if err = uc.userRepo.SetPasswordResetTOTP(ctx, user.Id, secure.GenerateHash(code)); err != nil {
		logger.Error(err)
		return shared.ErrInsertRecord
	}
	if err = uc.userRepo.UpdateLastPassReset(ctx, user.Id); err != nil {
		logger.Error(err)
		return shared.ErrUpdateRecord
	}

	msg := []byte(fmt.Sprintf(`Код для смены пароля: %s. Действителен в течение 5 минут`, code))
	return uc.mailServer.Send(user.Email, msg)
}

func (uc *userController) ValidateResetPasswordCode(ctx context.Context, email, code string) error {
	user, err := uc.userRepo.FindOneByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return shared.ErrNoSuchUser
		}
		return shared.ErrFindRecord
	}

	secret, err := uc.userRepo.GetPasswordResetTOTP(ctx, user.Id)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return shared.ErrNoTOTP
		}
		logger.Error(err)
		return shared.ErrFindRecord
	}

	if secret != secure.GenerateHash(code) {
		return shared.ErrInvalidTOTP
	}
	if err = uc.userRepo.DeletePasswordResetTOTP(ctx, user.Id); err != nil {
		return shared.ErrDeleteRecord
	}

	return uc.userRepo.SetResetPassAvailable(ctx, user.Id)
}

func (uc *userController) ResetPassword(ctx context.Context, email, password string) error {
	user, err := uc.userRepo.FindOneByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return shared.ErrNoSuchUser
		}
		return shared.ErrFindRecord
	}

	if err = uc.userRepo.GetDelResetPassAvailable(ctx, user.Id); err != nil {
		if errors.Is(err, redis.Nil) {
			return shared.ErrResetPassTimeout
		}
		return shared.ErrFindRecord
	}

	hashedPassword, err := secure.GetPasswordHash(password)
	if err != nil {
		return shared.ErrPasswordTooLong
	}

	return uc.userRepo.UpdatePassword(ctx, user.Id, hashedPassword)
}
