package handler

import (
	"purple/internal/shared"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"github.com/yogenyslav/logger"
)

type session struct {
	id uuid.UUID
}

type SessionWsHandler struct {
	sessions   map[*websocket.Conn]*session
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
}

func NewSessionWsHandler() *SessionWsHandler {
	return &SessionWsHandler{
		sessions:   make(map[*websocket.Conn]*session),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

func (h *SessionWsHandler) Serve() {
	for {
		select {
		case conn := <-h.register:
			id, err := uuid.Parse(conn.Params("id"))
			if err != nil {
				logger.Errorf("%v: %v", shared.ErrSessionIdInvalid, err)
				err = conn.WriteMessage(websocket.TextMessage, []byte(shared.ErrSessionIdInvalid.Error()))
				if err != nil {
					logger.Errorf("failed to write message: %v", err)
				}
				break
			}

			h.sessions[conn] = &session{
				id: id,
			}
			logger.Infof("started session with id %s", id.String())

		case conn := <-h.unregister:
			delete(h.sessions, conn)
			idRaw := conn.Params("id")
			id, err := uuid.Parse(idRaw)
			if err != nil {
				logger.Errorf("failed to delete session: %v :%v", shared.ErrSessionIdInvalid, err)
				break
			}
			logger.Infof("closed session with id %s", id.String())
		}
	}
}

func (h *SessionWsHandler) Chat(c *websocket.Conn) {
	var (
		mt  int
		msg []byte
		err error
	)

	defer func() {
		h.unregister <- c
		if err = c.Close(); err != nil {
			logger.Errorf("failed to close connection: %v", err)
		}
	}()

	h.register <- c

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
		if err = c.WriteMessage(mt, []byte("pupupu")); err != nil {
			logger.Error(err)
		}
	}
}
