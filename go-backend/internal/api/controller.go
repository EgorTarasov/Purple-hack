package api

import (
	"context"

	"purple/internal/api/domain"

	"github.com/google/uuid"
)

type SessionController interface {
	InsertOne(ctx context.Context, id uuid.UUID) error
	FindOneById(ctx context.Context, id uuid.UUID) (domain.Session, error)
	List(ctx context.Context, userId int64) ([]domain.Session, error)
}

type QueryController interface {
	InsertOne(ctx context.Context, params domain.QueryCreate) (domain.Query, error)
	FindMany(ctx context.Context, sessionId uuid.UUID) ([]domain.Query, error)
}

type ResponseController interface {
	InsertOne(ctx context.Context, params domain.ResponseCreate) (domain.Response, error)
	Respond(ctx context.Context, params domain.Query, sessionId uuid.UUID) (domain.ResponseCreate, error)
	RespondStream(
		ctx context.Context, params domain.Query, sessionId uuid.UUID,
		ctxCh, bodyCh chan<- string,
	) (domain.ResponseCreate, error)
	FindMany(ctx context.Context, sessionId uuid.UUID) ([]domain.Response, error)
}

type UserController interface {
	SaveSession(ctx context.Context, userId int64, sessionId uuid.UUID) error
}
