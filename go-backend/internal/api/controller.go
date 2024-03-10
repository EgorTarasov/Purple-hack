package api

import (
	"context"

	"purple/internal/api/domain"

	"github.com/google/uuid"
)

type SessionController interface {
	InsertOne(ctx context.Context, id uuid.UUID) error
	FindOne(ctx context.Context, id uuid.UUID) (domain.Session, error)
}

type QueryController interface {
	InsertOne(ctx context.Context, params domain.QueryCreate) (domain.Query, error)
}

type ResponseController interface {
	InsertOne(ctx context.Context, params domain.ResponseCreate) (domain.Response, error)
	Respond(ctx context.Context, params domain.Query, sessionId uuid.UUID) (domain.ResponseCreate, error)
}
