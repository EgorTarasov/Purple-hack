package handler

import (
	"context"
	"strconv"

	"purple/internal/api"
	"purple/internal/api/domain"
	"purple/internal/server/response"
	"purple/internal/shared"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"github.com/yogenyslav/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		mt             int
		msg            []byte
		err            error
		responseStatus string
		respCreate     domain.ResponseCreate
		queryCreate    domain.QueryCreate
		query          domain.Query
		resp           domain.Response
		respCtx        string
		respBody       string
		sessionId      uuid.UUID
		userId         int64
		data           []byte
		ctx            = context.Background()
		ctxCh          = make(chan string, 1)
		bodyCh         = make(chan string)
		errCh          = make(chan error, 1)
		cancelCh       = make(chan int)
		respCh         = make(chan domain.Response, 1)
	)

	defer func() {
		if err = c.Close(); err != nil {
			logger.Errorf("failed to close connection: %v", err)
		}
	}()

	logger.Infof("new conn from %s", c.LocalAddr())

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

		if string(msg) == response.StreamCancel {
			cancelCh <- 1
			continue
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
				responseStatus = "failed to save session"
				logger.Errorf("%s: %v", responseStatus, err)
				if err = c.WriteMessage(websocket.CloseMessage, []byte(responseStatus)); err != nil {
					logger.Error(err)
					break
				}
				break
			}
		}

		if err = h.sessionController.InsertOne(ctx, sessionId); err != nil {
			responseStatus = "failed to create session"
			logger.Errorf("%s: %v", responseStatus, err)
			if err = c.WriteMessage(websocket.CloseMessage, []byte(responseStatus)); err != nil {
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
			responseStatus = "failed to create query"
			logger.Errorf("%s: %v", responseStatus, err)
			if err = c.WriteMessage(websocket.CloseMessage, []byte(responseStatus)); err != nil {
				logger.Error(err)
				break
			}
			break
		}

		go func() {
			respCreate, err = h.responseController.RespondStream(ctx, query, sessionId, ctxCh, bodyCh, cancelCh)
			if err != nil {
				logger.Info("failed to get response from searchEngine")
				logger.Error(err)
				errCh <- err
				return
			}

			resp, err = h.responseController.InsertOne(ctx, respCreate)
			if err != nil {
				logger.Info("failed to create response")
				logger.Error(err)
				errCh <- err
				return
			}

			respCh <- resp
		}()

		go func() {
		ReceiveStream:
			for {
				select {
				case streamErr := <-errCh:
					data = []byte(response.StreamError)

					grpcErr := status.Convert(streamErr)
					if grpcErr.Code() == codes.Canceled {
						logger.Info("stream canceled")
						data = []byte(response.StreamFinished)
					}

					logger.Infof("sending error: %v with data %s", streamErr, data)
					if err = c.WriteMessage(websocket.TextMessage, data); err != nil {
						responseStatus = "failed to write grpc stream error"
						logger.Errorf("%s: %v", responseStatus, err)
						if err = c.WriteMessage(websocket.CloseMessage, []byte(responseStatus)); err != nil {
							logger.Error(err)
						}
					}
					return
				case resp = <-respCh:
					logger.Infof("final resp: %v", resp)
					break ReceiveStream
				case respCtx = <-ctxCh:
					_ = respCtx
					//if err = c.WriteMessage(websocket.TextMessage, []byte(respCtx)); err != nil {
					//	responseStatus = "failed to write message"
					//	logger.Errorf("%s: %v", responseStatus, err)
					//	if err = c.WriteMessage(websocket.CloseMessage, []byte(responseStatus)); err != nil {
					//		logger.Error(err)
					//	}
					//	break ReceiveStream
					//}
				case respBody = <-bodyCh:
					if err = c.WriteMessage(websocket.TextMessage, []byte(respBody)); err != nil {
						responseStatus = "failed to write message"
						logger.Errorf("%s: %v", responseStatus, err)
						if err = c.WriteMessage(websocket.CloseMessage, []byte(responseStatus)); err != nil {
							logger.Error(err)
						}
						break ReceiveStream
					}
				}
			}
			logger.Info("stream finished")
			if err = c.WriteMessage(websocket.TextMessage, []byte(response.StreamFinished)); err != nil {
				responseStatus = "failed to write message"
				logger.Errorf("%s: %v", responseStatus, err)
				if err = c.WriteMessage(websocket.CloseMessage, []byte(responseStatus)); err != nil {
					logger.Error(err)
				}
			}
		}()
	}
}
