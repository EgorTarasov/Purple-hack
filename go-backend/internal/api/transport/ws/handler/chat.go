package handler

import (
	"context"
	"encoding/json"

	"purple/internal/api"
	"purple/internal/api/domain"
	"purple/internal/shared"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"github.com/yogenyslav/logger"
)

type ChatHandler struct {
	sessionController  api.SessionController
	queryController    api.QueryController
	responseController api.ResponseController
}

func NewChatHandler(sc api.SessionController, qc api.QueryController, rc api.ResponseController) *ChatHandler {
	return &ChatHandler{
		sessionController:  sc,
		queryController:    qc,
		responseController: rc,
	}
}

func (h *ChatHandler) Chat(c *websocket.Conn) {
	var (
		mt         int
		msg        []byte
		err        error
		status     string
		respCreate domain.ResponseCreate
		query      domain.Query
		resp       domain.Response
		data       []byte
		sessionId  uuid.UUID
		ctx        = context.Background()
	)

	defer func() {
		if err = c.Close(); err != nil {
			logger.Errorf("failed to close connection: %v", err)
		}
	}()

	for {
		mt, msg, err = c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Errorf("ws read error: %v", err)
				return
			}
			logger.Error(err)
			return
		}
		logger.Infof("got message %s", msg)

		if mt != websocket.TextMessage {
			logger.Error(shared.ErrUnexpectedMessageType)
			if err = c.WriteMessage(websocket.CloseMessage, []byte(shared.ErrUnexpectedMessageType.Error())); err != nil {
				logger.Error(err)
				return
			}
		}

		sessionId, err = uuid.Parse(c.Params("id"))
		if err != nil {
			logger.Errorf("%v: %v", shared.ErrSessionIdInvalid, err)
			if err = c.WriteMessage(websocket.CloseMessage, []byte(shared.ErrSessionIdInvalid.Error())); err != nil {
				logger.Error(err)
				return
			}
			return
		}

		if err = h.sessionController.InsertOne(ctx, sessionId); err != nil {
			status = "failed to create session"
			logger.Errorf("%s: %v", status, err)
			if err = c.WriteMessage(websocket.CloseMessage, []byte(status)); err != nil {
				logger.Error(err)
				return
			}
			return
		}

		queryCreate := domain.QueryCreate{
			SessionId: sessionId,
			Model:     c.Query("model", "llama"),
			Body:      string(msg),
			UserAgent: c.Headers("User-Agent"),
		}
		query, err = h.queryController.InsertOne(ctx, queryCreate)
		if err != nil {
			status = "failed to create query"
			logger.Errorf("%s: %v", status, err)
			if err = c.WriteMessage(websocket.CloseMessage, []byte(status)); err != nil {
				logger.Error(err)
				return
			}
			return
		}

		respCreate, err = h.responseController.Respond(ctx, query, sessionId)
		if err != nil {
			status = "failed to get response from searchEngine"
			logger.Errorf("%s: %v", status, err)
			if err = c.WriteMessage(websocket.CloseMessage, []byte(status)); err != nil {
				logger.Error(err)
				return
			}
			return
		}

		resp, err = h.responseController.InsertOne(ctx, respCreate)
		if err != nil {
			status = "failed to create response"
			logger.Errorf("%s: %v", status, err)
			if err = c.WriteMessage(websocket.CloseMessage, []byte(status)); err != nil {
				logger.Error(err)
				return
			}
			return
		}

		logger.Debugf("final resp: %v", resp)
		data, err = json.Marshal(resp)
		if err != nil {
			status = "failed to marshal resp"
			logger.Errorf("%s: %v", status, err)
			if err = c.WriteMessage(websocket.CloseMessage, []byte(status)); err != nil {
				logger.Error(err)
				return
			}
			return
		}

		if err = c.WriteMessage(websocket.TextMessage, data); err != nil {
			status = "failed to write message"
			logger.Errorf("%s: %v", status, err)
			if err = c.WriteMessage(websocket.CloseMessage, []byte(status)); err != nil {
				logger.Error(err)
				return
			}
			return
		}
	}
}
