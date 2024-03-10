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
	select id, query_ids, response_ids, created_at
	from session
	where id=$1;
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
