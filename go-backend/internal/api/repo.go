package api

import (
	"context"

	"purple/internal/api/data"

	"github.com/google/uuid"
)

type SessionRepo interface {
	InsertOne(ctx context.Context, params data.Session) error
	FindOne(ctx context.Context, id uuid.UUID) (data.Session, error)
}

type QueryRepo interface {
	InsertOne(ctx context.Context, params data.Query) (int64, error)
	FindMany(ctx context.Context, sessionId uuid.UUID) ([]data.Query, error)
}

type ResponseRepo interface {
	InsertOne(ctx context.Context, params data.Response) (int64, error)
	FindMany(ctx context.Context, sessionId uuid.UUID) ([]data.Response, error)
}
