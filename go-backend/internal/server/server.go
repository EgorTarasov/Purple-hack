package server

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"purple/config"
	"purple/internal/api/transport/ws"
	wsHandler "purple/internal/api/transport/ws/handler"
	"purple/internal/server/response"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/swagger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	loggermw "github.com/gofiber/fiber/v2/middleware/logger"
	recovermw "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yogenyslav/logger"
	"github.com/yogenyslav/storage/postgres"
)

type Server struct {
	app *fiber.App
	cfg *config.Config
	pg  *pgxpool.Pool
}

func New(cfg *config.Config) *Server {
	app := fiber.New(fiber.Config{
		ServerHeader: "Fiber",
		BodyLimit:    500 * 1024 * 1024,
		ErrorHandler: response.ErrorHandler,
	})

	pg := postgres.MustNew(&cfg.Postgres, 20)

	app.Use(loggermw.New())
	app.Use(recovermw.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Server.CorsOrigins,
		AllowMethods:     "*",
		AllowHeaders:     "*",
		AllowCredentials: true,
	}))

	app.Get("/swagger/*", swagger.HandlerDefault)

	return &Server{
		app: app,
		cfg: cfg,
		pg:  pg,
	}
}

func (s *Server) Run() {
	defer s.pg.Close()

	wsConfig := websocket.Config{
		Origins: strings.Split(s.cfg.Server.CorsOrigins, ","),
		RecoverHandler: func(conn *websocket.Conn) {
			if err := recover(); err != nil {
				err = conn.WriteJSON(fiber.Map{"error": fmt.Sprintf("internal error: %v", err)})
				if err != nil {
					logger.Errorf("failed to handle ws error: %v", err)
				}
			}
		},
	}
	sessionWsHandler := wsHandler.NewSessionWsHandler()
	ws.SetupSessionSocket(s.app, sessionWsHandler, &wsConfig)

	go sessionWsHandler.Serve()
	go s.listen()
	s.gracefulShutdown()
}

func (s *Server) listen() {
	addr := fmt.Sprintf(":%d", s.cfg.Server.Port)
	logger.Infof("starting Server %s", addr)
	if err := s.app.Listen(addr); err != nil {
		logger.Fatalf("error has occurred while listening on %s: %v", addr, err)
	}
}

func (s *Server) gracefulShutdown() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	_ = <-ch

	s.pg.Close()
	logger.Info("shutting down the server")
}
