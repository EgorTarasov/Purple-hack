package repo

import (
	"context"

	"purple/internal/api/data"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yogenyslav/storage/postgres"
)

type QueryRepo struct {
	pg *pgxpool.Pool
}

const insertOneQuery = `
	insert into query(fk_session_id, body, model, user_agent, created_at)
	values ($1, $2, $3, $4, $5)
	returning id;
`

const findManyQuery = `
	select id, fk_session_id, body, model, user_agent, created_at
	from query
	where fk_session_id=$1;
`

func NewQueryRepo(pg *pgxpool.Pool) *QueryRepo {
	return &QueryRepo{
		pg: pg,
	}
}

func (r *QueryRepo) InsertOne(ctx context.Context, params data.Query) (int64, error) {
	return postgres.QueryPrimitive[int64](
		ctx,
		r.pg,
		insertOneQuery,
		params.SessionId,
		params.Body,
		params.Model,
		params.UserAgent,
		params.CreatedAt,
	)
}

func (r *QueryRepo) FindMany(ctx context.Context, sessionId uuid.UUID) ([]data.Query, error) {
	return postgres.QueryStructSlice[data.Query](ctx, r.pg, findManyQuery, sessionId)
}
