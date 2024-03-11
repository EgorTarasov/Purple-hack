package controller

import (
	"context"

	"purple/internal/api"
	"purple/internal/api/data"
	"purple/internal/api/domain"
	"purple/internal/shared"
	"purple/pkg"

	"github.com/google/uuid"
	"github.com/yogenyslav/logger"
)

type SessionController struct {
	sessionRepo  api.SessionRepo
	queryRepo    api.QueryRepo
	responseRepo api.ResponseRepo
}

func NewSessionController(sr api.SessionRepo) *SessionController {
	return &SessionController{
		sessionRepo: sr,
	}
}

func (sc *SessionController) InsertOne(ctx context.Context, id uuid.UUID) error {
	session := data.Session{
		Id:        id,
		CreatedAt: pkg.GetLocalTime(),
	}
	if err := sc.sessionRepo.InsertOne(ctx, session); err != nil {
		if pkg.CheckErrDuplicateKey(err) {
			// FIXME: handle duplicating uuids
			logger.Debugf("duplicating id: %s", id.String())
			return nil
		}
		logger.Error(err)
		return shared.ErrInsertRecord
	}
	return nil
}

func (sc *SessionController) FindOneById(ctx context.Context, id uuid.UUID) (domain.Session, error) {
	var (
		session domain.Session
		err     error
	)

	sessionDb, err := sc.sessionRepo.FindOne(ctx, id)
	if err != nil {
		logger.Error(err)
		return session, shared.ErrFindRecord
	}

	session = domain.Session{
		Id:        sessionDb.Id,
		CreatedAt: sessionDb.CreatedAt,
	}

	queries, responses, err := sc.listQueriesResponses(ctx, id)
	if err != nil {
		return session, err
	}

	session.Queries = queries
	session.Responses = responses
	return session, nil
}

func (sc *SessionController) List(ctx context.Context, userId int64) ([]domain.Session, error) {
	sessionsDb, err := sc.sessionRepo.List(ctx, userId)
	if err != nil {
		logger.Error(err)
		return nil, shared.ErrFindRecord
	}

	sessions := make([]domain.Session, 0, len(sessionsDb))
	for _, session := range sessionsDb {
		queries, responses, err := sc.listQueriesResponses(ctx, session.Id)
		if err != nil {
			return nil, err
		}

		sessions = append(sessions, domain.Session{
			Id:        session.Id,
			Queries:   queries,
			Responses: responses,
			CreatedAt: session.CreatedAt,
		})
	}

	return sessions, nil
}

func (sc *SessionController) listQueriesResponses(ctx context.Context, sessionId uuid.UUID) ([]domain.Query, []domain.Response, error) {
	queriesDb, err := sc.queryRepo.FindMany(ctx, sessionId)
	if err != nil {
		logger.Error(err)
		return nil, nil, shared.ErrFindRecord
	}
	queries := make([]domain.Query, 0, len(queriesDb))
	for _, query := range queriesDb {
		queries = append(queries, domain.Query{
			Id:        query.Id,
			Model:     query.Model,
			Body:      query.Body,
			CreatedAt: query.CreatedAt,
		})
	}

	responsesDb, err := sc.responseRepo.FindMany(ctx, sessionId)
	if err != nil {
		logger.Error(err)
		return nil, nil, shared.ErrFindRecord
	}
	responses := make([]domain.Response, 0, len(responsesDb))
	for _, response := range responsesDb {
		responses = append(responses, domain.Response{
			Id:        response.Id,
			Body:      response.Body,
			Context:   response.Context,
			CreatedAt: response.CreatedAt,
		})
	}

	return queries, responses, nil
}
