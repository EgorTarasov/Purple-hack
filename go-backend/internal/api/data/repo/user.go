package repo

import (
	"context"

	"purple/internal/api/data"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yogenyslav/storage/postgres"
)

type UserRepo struct {
	pg *pgxpool.Pool
}

const findOneUser = `
	select id, email, password, created_at
	from "user"
	where email=$1;
`

const saveSession = `
	insert into users_sessions(fk_user_id, fk_session_id)
	values ($1, $2);
`

func NewUserRepo(pg *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		pg: pg,
	}
}

func (r *UserRepo) FindOneByEmail(ctx context.Context, email string) (data.User, error) {
	return postgres.QueryStruct[data.User](ctx, r.pg, findOneUser, email)
}

func (r *UserRepo) SaveSession(ctx context.Context, userId int64, sessionId uuid.UUID) error {
	_, err := r.pg.Query(ctx, saveSession, userId, sessionId)
	return err
}
