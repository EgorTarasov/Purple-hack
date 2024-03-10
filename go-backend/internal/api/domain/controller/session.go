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

func NewSessionController(
	sr api.SessionRepo,
	// qr api.QueryRepo,
	// rr api.ResponseRepo,
) *SessionController {
	return &SessionController{
		sessionRepo: sr,
		//queryRepo:    qr,
		//responseRepo: rr,
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
			logger.Debugf("dubplicating id: %s", id.String())
			return nil
		}
		logger.Error(err)
		return shared.ErrInsertRecord
	}
	return nil
}

func (sc *SessionController) FindOne(ctx context.Context, id uuid.UUID) (domain.Session, error) {
	var (
		session domain.Session
		err     error
	)

	sessionDb, err := sc.sessionRepo.FindOne(ctx, id)
	if err != nil {
		logger.Error(err)
		return session, shared.ErrFindRecord
	}

	//queriesDb, err := sc.queryRepo.FindMany(ctx, id)
	//if err != nil {
	//	logger.Error(err)
	//	return session, shared.ErrFindRecord
	//}
	//
	//queries := make([]domain.Query, len(queriesDb))
	//for _, query := range queriesDb {
	//	queries = append(queries, domain.Query{
	//		Id:        query.Id,
	//		Body:      query.Body,
	//		CreatedAt: query.CreatedAt,
	//	})
	//}
	//
	//responsesDb, err := sc.responseRepo.FindMany(ctx, id)
	//if err != nil {
	//	logger.Error(err)
	//	return session, shared.ErrFindRecord
	//}
	//
	//responses := make([]domain.Response, len(responsesDb))
	//for _, response := range responsesDb {
	//	responses = append(responses, domain.Response{
	//		Id:        response.Id,
	//		Body:      response.Body,
	//		Context:   response.Context,
	//		CreatedAt: response.CreatedAt,
	//	})
	//}

	session = domain.Session{
		Id:          sessionDb.Id,
		QueryIds:    sessionDb.QueryIds,
		ResponseIds: sessionDb.ResponseIds,
		CreatedAt:   sessionDb.CreatedAt,
	}
	return session, nil
}
