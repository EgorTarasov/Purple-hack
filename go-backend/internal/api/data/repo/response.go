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

const updateSessionResponse = `
	update session 
	set response_ids = array_append(response_ids, $1)
	where session.id = $2;
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
	responseId, err := postgres.QueryPrimitive[int64](
		ctx,
		r.pg,
		insertOneResponse,
		params.SessionId,
		params.QueryId,
		params.Body,
		params.Context,
		params.CreatedAt,
	)

	if err != nil {
		return 0, err
	}

	if _, err = r.pg.Exec(ctx, updateSessionResponse, responseId, params.SessionId); err != nil {
		return 0, err
	}
	return responseId, nil
}

func (r *ResponseRepo) FindMany(ctx context.Context, sessionId uuid.UUID) ([]data.Response, error) {
	return postgres.QueryStruct[[]data.Response](ctx, r.pg, findManyResponse, sessionId)
}
