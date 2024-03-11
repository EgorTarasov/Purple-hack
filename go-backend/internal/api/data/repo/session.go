package repo

import (
	"context"

	"purple/internal/api/data"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yogenyslav/storage/postgres"
)

type SessionRepo struct {
	pg *pgxpool.Pool
}

const insertOneSession = `
	insert into session(id, created_at)
	values ($1, $2)
	returning id;
`

const findOneSession = `
	select id, created_at
	from session
	where id=$1;
`

const listSessions = `
	select id, created_at
	from session
	where id in (
	    select fk_session_id
	    from users_sessions
	    where fk_user_id=$1
	);
`

func NewSessionRepo(pg *pgxpool.Pool) *SessionRepo {
	return &SessionRepo{
		pg: pg,
	}
}

func (r *SessionRepo) InsertOne(ctx context.Context, params data.Session) error {
	_, err := postgres.QueryPrimitive[uuid.UUID](ctx, r.pg, insertOneSession, params.Id, params.CreatedAt)
	return err
}

func (r *SessionRepo) FindOne(ctx context.Context, id uuid.UUID) (data.Session, error) {
	return postgres.QueryStruct[data.Session](ctx, r.pg, findOneSession, id)
}

func (r *SessionRepo) List(ctx context.Context, userId int64) ([]data.Session, error) {
	return postgres.QueryStructSlice[data.Session](ctx, r.pg, listSessions, userId)
}
