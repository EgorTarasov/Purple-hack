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

type QueryController struct {
	repo api.QueryRepo
}

func NewQueryController(qr api.QueryRepo) *QueryController {
	return &QueryController{
		repo: qr,
	}
}

func (qc *QueryController) InsertOne(ctx context.Context, params domain.QueryCreate) (domain.Query, error) {
	var (
		result domain.Query
		query  data.Query
	)

	query = data.Query{
		SessionId: params.SessionId,
		Body:      params.Body,
		Model:     params.Model,
		UserAgent: params.UserAgent,
		CreatedAt: pkg.GetLocalTime(),
	}
	queryId, err := qc.repo.InsertOne(ctx, query)
	if err != nil {
		logger.Error(err)
		return result, shared.ErrInsertRecord
	}

	result = domain.Query{
		Id:        queryId,
		Model:     query.Model,
		Body:      query.Body,
		CreatedAt: query.CreatedAt,
	}
	return result, nil
}

func (qc *QueryController) FindMany(ctx context.Context, sessionId uuid.UUID) ([]domain.Query, error) {
	queriesDb, err := qc.repo.FindMany(ctx, sessionId)
	if err != nil {
		logger.Error(err)
		return nil, shared.ErrFindRecord
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

	return queries, nil
}
