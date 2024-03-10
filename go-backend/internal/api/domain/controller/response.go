package controller

import (
	"context"

	"purple/internal/api"
	"purple/internal/api/data"
	"purple/internal/api/domain"
	"purple/internal/shared"
	"purple/pkg"
	protos "purple/proto"

	"github.com/google/uuid"
	"github.com/yogenyslav/logger"
	"google.golang.org/grpc"
)

type ResponseController struct {
	repo     api.ResponseRepo
	seClient protos.SearchEngineClient
}

func NewResponseController(rr api.ResponseRepo, seConn *grpc.ClientConn) *ResponseController {
	return &ResponseController{
		repo:     rr,
		seClient: protos.NewSearchEngineClient(seConn),
	}
}

func (rc *ResponseController) InsertOne(ctx context.Context, params domain.ResponseCreate) (domain.Response, error) {
	var (
		resp   domain.Response
		respDb data.Response
	)

	respDb = data.Response{
		SessionId: params.SessionId,
		QueryId:   params.QueryId,
		Body:      params.Body,
		Context:   params.Context,
		CreatedAt: pkg.GetLocalTime(),
	}
	respId, err := rc.repo.InsertOne(ctx, respDb)
	if err != nil {
		logger.Error(err)
		return resp, shared.ErrInsertRecord
	}

	resp = domain.Response{
		Id:        respId,
		Body:      respDb.Body,
		Context:   respDb.Context,
		CreatedAt: respDb.CreatedAt,
	}
	return resp, nil
}

func (rc *ResponseController) Respond(ctx context.Context, params domain.Query, sessionId uuid.UUID) (domain.ResponseCreate, error) {
	var (
		resp     domain.ResponseCreate
		respGrpc *protos.Response
		err      error
	)

	in := protos.Query{
		Body:  params.Body,
		Model: params.Model,
	}
	respGrpc, err = rc.seClient.Respond(ctx, &in)
	if err != nil {
		return resp, err
	}

	respContext := make(map[string][]string)
	for id, obj := range respGrpc.GetContext() {
		respContext[id] = obj.GetValue()
	}

	resp = domain.ResponseCreate{
		SessionId: sessionId,
		QueryId:   params.Id,
		Body:      respGrpc.GetBody(),
		Context:   respContext,
		Model:     params.Model,
	}
	return resp, nil
}
