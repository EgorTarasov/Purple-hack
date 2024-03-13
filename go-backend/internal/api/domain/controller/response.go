package controller

import (
	"context"
	"io"

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

// Respond deprecated
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

	resp = domain.ResponseCreate{
		SessionId: sessionId,
		QueryId:   params.Id,
		Body:      respGrpc.GetBody(),
		Context:   respGrpc.GetContext(),
		Model:     params.Model,
	}
	return resp, nil
}

func (rc *ResponseController) RespondStream(
	ctx context.Context, params domain.Query, sessionId uuid.UUID,
	ctxCh, bodyCh chan<- string, cancelCh <-chan int,
) (domain.ResponseCreate, error) {
	var (
		tmp         string
		resp        domain.ResponseCreate
		err         error
		stream      protos.SearchEngine_RespondStreamClient
		respGrpc    *protos.Response
		respBody    string
		respContext string
	)

	withCancel, cancel := context.WithCancel(ctx)
	defer cancel()

	in := protos.Query{
		Body:  params.Body,
		Model: params.Model,
	}
	stream, err = rc.seClient.RespondStream(withCancel, &in)
	if err != nil {
		return resp, err
	}

reading:
	for {
		select {
		case <-cancelCh:
			cancel()
		default:
			respGrpc, err = stream.Recv()
			if err == io.EOF {
				break reading
			}
			if err != nil {
				return resp, err
			}
			logger.Debugf("got from stream: %s", respGrpc.String())
			if tmp = respGrpc.GetBody(); tmp != "" {
				respBody += tmp
				bodyCh <- tmp
			}
			if tmp = respGrpc.GetContext(); tmp != "" {
				respContext += tmp
				ctxCh <- tmp
			}
		}
	}

	resp = domain.ResponseCreate{
		SessionId: sessionId,
		QueryId:   params.Id,
		Body:      respBody,
		Context:   respContext,
		Model:     params.Model,
	}
	return resp, nil
}

func (rc *ResponseController) FindMany(ctx context.Context, sessionId uuid.UUID) ([]domain.Response, error) {
	responsesDb, err := rc.repo.FindMany(ctx, sessionId)
	if err != nil {
		logger.Error(err)
		return nil, shared.ErrFindRecord
	}

	responses := make([]domain.Response, 0, len(responsesDb))
	for _, resp := range responsesDb {
		responses = append(responses, domain.Response{
			Id:        resp.Id,
			Body:      resp.Body,
			Context:   resp.Context,
			CreatedAt: resp.CreatedAt,
		})
	}

	return responses, nil
}
