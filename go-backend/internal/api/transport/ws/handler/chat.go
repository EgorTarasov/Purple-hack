package handler

import (
	"context"
	"encoding/json"
	"strconv"

	"purple/internal/api"
	"purple/internal/api/domain"
	"purple/internal/server/response"
	"purple/internal/shared"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"github.com/yogenyslav/logger"
)

type ChatHandler struct {
	sessionController  api.SessionController
	queryController    api.QueryController
	responseController api.ResponseController
	userController     api.UserController
}

func NewChatHandler(sc api.SessionController, qc api.QueryController, rc api.ResponseController, uc api.UserController) *ChatHandler {
	return &ChatHandler{
		sessionController:  sc,
		queryController:    qc,
		responseController: rc,
		userController:     uc,
	}
}

func (h *ChatHandler) Chat(c *websocket.Conn) {
	var (
		mt          int
		msg         []byte
		err         error
		status      string
		respCreate  domain.ResponseCreate
		queryCreate domain.QueryCreate
		query       domain.Query
		resp        domain.Response
		respCtx     string
		respBody    string
		data        []byte
		sessionId   uuid.UUID
		userId      int64
		ctx         = context.Background()
		ctxCh       = make(chan string, 1)
		bodyCh      = make(chan string)
		errCh       = make(chan error, 1)
		respCh      = make(chan domain.Response, 1)
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
				break
			}
			logger.Error(err)
			break
		}
		logger.Infof("got message %s", msg)

		if mt != websocket.TextMessage {
			logger.Error(shared.ErrUnexpectedMessageType)
			if err = c.WriteMessage(websocket.CloseMessage, []byte(shared.ErrUnexpectedMessageType.Error())); err != nil {
				logger.Error(err)
				break
			}
		}

		sessionId, err = uuid.Parse(c.Params("id"))
		if err != nil {
			logger.Errorf("%v: %v", shared.ErrSessionIdInvalid, err)
			if err = c.WriteMessage(websocket.CloseMessage, []byte(shared.ErrSessionIdInvalid.Error())); err != nil {
				logger.Error(err)
				break
			}
			break
		}

		userId, err = strconv.ParseInt(c.Cookies("auth"), 10, 64)
		if err == nil {
			if err = h.userController.SaveSession(ctx, userId, sessionId); err != nil {
				status = "failed to save session"
				logger.Errorf("%s: %v", status, err)
				if err = c.WriteMessage(websocket.CloseMessage, []byte(status)); err != nil {
					logger.Error(err)
					break
				}
				break
			}
		}

		if err = h.sessionController.InsertOne(ctx, sessionId); err != nil {
			status = "failed to create session"
			logger.Errorf("%s: %v", status, err)
			if err = c.WriteMessage(websocket.CloseMessage, []byte(status)); err != nil {
				logger.Error(err)
				break
			}
			break
		}

		queryCreate = domain.QueryCreate{
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
				break
			}
			break
		}

		// TODO: replace context to WithCancel in order to cancel grpc streaming
		go func() {
			respCreate, err = h.responseController.RespondStream(ctx, query, sessionId, ctxCh, bodyCh)
			if err != nil {
				status = "failed to get response from searchEngine"
				logger.Errorf("%s: %v", status, err)
				errCh <- err
				return
			}

			resp, err = h.responseController.InsertOne(ctx, respCreate)
			if err != nil {
				status = "failed to create response"
				logger.Errorf("%s: %v", status, err)
				errCh <- err
				return
			}

			respCh <- resp
		}()

	ReceiveStream:
		for {
			select {
			case err = <-errCh:
				logger.Errorf("sending error: %v", err)
				if err = c.WriteMessage(websocket.TextMessage, []byte(response.StreamError)); err != nil {
					status = "failed to write grpc stream error"
					logger.Errorf("%s: %v", status, err)
					if err = c.WriteMessage(websocket.CloseMessage, []byte(status)); err != nil {
						logger.Error(err)
					}
					break ReceiveStream
				}
			case resp = <-respCh:
				logger.Infof("final resp: %v", resp)
				data, err = json.Marshal(resp)
				if err != nil {
					status = "failed to marshal resp"
					logger.Errorf("%s: %v", status, err)
					if err = c.WriteMessage(websocket.CloseMessage, []byte(status)); err != nil {
						logger.Error(err)
					}
					break ReceiveStream
				}

				if err = c.WriteMessage(websocket.TextMessage, data); err != nil {
					status = "failed to write message"
					logger.Errorf("%s: %v", status, err)
					if err = c.WriteMessage(websocket.CloseMessage, []byte(status)); err != nil {
						logger.Error(err)
					}
				}
				break ReceiveStream
			case respCtx = <-ctxCh:
				_ = respCtx
				//if err = c.WriteMessage(websocket.TextMessage, []byte(respCtx)); err != nil {
				//	status = "failed to write message"
				//	logger.Errorf("%s: %v", status, err)
				//	if err = c.WriteMessage(websocket.CloseMessage, []byte(status)); err != nil {
				//		logger.Error(err)
				//	}
				//	break ReceiveStream
				//}
			case respBody = <-bodyCh:
				if err = c.WriteMessage(websocket.TextMessage, []byte(respBody)); err != nil {
					status = "failed to write message"
					logger.Errorf("%s: %v", status, err)
					if err = c.WriteMessage(websocket.CloseMessage, []byte(status)); err != nil {
						logger.Error(err)
					}
					break ReceiveStream
				}
			}
		}
		logger.Info("stream finished")
		if err = c.WriteMessage(websocket.TextMessage, []byte(response.StreamFinished)); err != nil {
			status = "failed to write message"
			logger.Errorf("%s: %v", status, err)
			if err = c.WriteMessage(websocket.CloseMessage, []byte(status)); err != nil {
				logger.Error(err)
			}
			break
		}
	}
}
