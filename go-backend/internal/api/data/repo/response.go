package repo

import (
	"context"

	"purple/internal/api/data"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yogenyslav/storage/postgres"
)

type ResponseRepo struct {
	pg *pgxpool.Pool
}

const insertOneResponse = `
	insert into response(fk_session_id, query_id, body, context, created_at)
	values ($1, $2, $3, $4, $5)
	returning id;
`

const findManyResponse = `
	select id, fk_session_id, query_id, body, context, created_at
	from response
	where fk_session_id=$1;
`

func NewResponseRepo(pg *pgxpool.Pool) *ResponseRepo {
	return &ResponseRepo{
		pg: pg,
	}
}

func (r *ResponseRepo) InsertOne(ctx context.Context, params data.Response) (int64, error) {
	return postgres.QueryPrimitive[int64](
		ctx,
		r.pg,
		insertOneResponse,
		params.SessionId,
		params.QueryId,
		params.Body,
		params.Context,
		params.CreatedAt,
	)
}

func (r *ResponseRepo) FindMany(ctx context.Context, sessionId uuid.UUID) ([]data.Response, error) {
	return postgres.QueryStructSlice[data.Response](ctx, r.pg, findManyResponse, sessionId)
}
