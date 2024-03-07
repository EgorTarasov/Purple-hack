package repo

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/yogenyslav/storage/postgres"
	"hack/internal/api"
	"hack/internal/api/data"
	"hack/internal/shared"
)

const (
	insertOne = `
		insert into users(
			email, 
			password, 
			first_name, 
			last_name
		)
		values ($1, $2, $3, $4)
		returning id;
	`

	findOneByEmail = `
		select id, 
		       email, 
		       password, 
		       first_name, 
		       last_name, 
		       last_pass_reset,
		       created_at, 
		       updated_at 
		from users
		where email = $1;
	`

	findOneById = `
		select id, 
		       email, 
		       password, 
		       first_name, 
		       last_name, 
		       last_pass_reset,
		       created_at, 
		       updated_at 
		from users
		where id = $1;
	`

	updateLastPassReset = `
		update users
		set last_pass_reset = current_timestamp,
		    updated_at = current_timestamp
		where id = $1;
	`

	updatePassword = `
		update users
		set password = $1,
		    updated_at = current_timestamp
		where id = $2;
	`
)

const (
	resetPasswordTotpKey      = "reset-pass-totp"
	resetPasswordAvailableKey = "reset-pass-available"
)

type userRepo struct {
	pg *pgxpool.Pool
	rc *redis.Client
}

func NewUserRepo(pg *pgxpool.Pool, rc *redis.Client) api.UserRepo {
	return &userRepo{
		pg: pg,
		rc: rc,
	}
}

func (r *userRepo) InsertOne(ctx context.Context, user data.User) (int64, error) {
	return postgres.QueryPrimitive[int64](
		ctx,
		r.pg,
		insertOne,
		user.Email,
		user.Password,
		user.FirstName,
		user.LastName,
	)
}

func (r *userRepo) FindOneByEmail(ctx context.Context, email string) (data.User, error) {
	return postgres.QueryStruct[data.User](ctx, r.pg, findOneByEmail, email)
}

func (r *userRepo) FindOneById(ctx context.Context, userId int64) (data.User, error) {
	return postgres.QueryStruct[data.User](ctx, r.pg, findOneById, userId)
}

func (r *userRepo) UpdateLastPassReset(ctx context.Context, userId int64) error {
	_, err := r.pg.Exec(ctx, updateLastPassReset, userId)
	return err
}

func (r *userRepo) UpdatePassword(ctx context.Context, userId int64, password string) error {
	_, err := r.pg.Exec(ctx, updatePassword, password, userId)
	return err
}

func (r *userRepo) SetPasswordResetTOTP(ctx context.Context, userId int64, secret string) error {
	return r.rc.Set(
		ctx,
		fmt.Sprintf("%s-%d", resetPasswordTotpKey, userId),
		secret,
		shared.PasswordResetExp,
	).Err()
}

func (r *userRepo) GetPasswordResetTOTP(ctx context.Context, userId int64) (string, error) {
	return r.rc.Get(ctx, fmt.Sprintf("%s-%d", resetPasswordTotpKey, userId)).Result()
}

func (r *userRepo) DeletePasswordResetTOTP(ctx context.Context, userId int64) error {
	return r.rc.Del(ctx, fmt.Sprintf("%s-%d", resetPasswordTotpKey, userId)).Err()
}

func (r *userRepo) SetResetPassAvailable(ctx context.Context, userId int64) error {
	return r.rc.Set(
		ctx,
		fmt.Sprintf("%s-%d", resetPasswordAvailableKey, userId),
		true,
		shared.PasswordResetExp,
	).Err()
}

func (r *userRepo) GetDelResetPassAvailable(ctx context.Context, userId int64) error {
	return r.rc.GetDel(ctx, fmt.Sprintf("%s-%d", resetPasswordAvailableKey, userId)).Err()
}
