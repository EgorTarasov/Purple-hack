package server

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"purple/config"
	"purple/internal/api/data/repo"
	"purple/internal/api/domain/controller"
	"purple/internal/api/transport/ws"
	wsHandler "purple/internal/api/transport/ws/handler"
	"purple/internal/server/response"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

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

	var grpcOpts []grpc.DialOption
	grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	seAddr := fmt.Sprintf("%s:%d", s.cfg.SE.Host, s.cfg.SE.Port)
	searchEngineConn, err := grpc.Dial(seAddr, grpcOpts...)
	if err != nil {
		logger.Panicf("failed to connect to searchEngine: %v", err)
	}
	defer func() {
		if err = searchEngineConn.Close(); err != nil {
			logger.Error(err)
		}
	}()

	wsConfig := websocket.Config{
		Origins: strings.Split(s.cfg.Server.CorsOrigins, ","),
		RecoverHandler: func(conn *websocket.Conn) {
			if e := recover(); e != nil {
				err = conn.WriteJSON(fiber.Map{"error": fmt.Sprintf("internal error: %v", err)})
				if err != nil {
					logger.Errorf("failed to handle ws error: %v", err)
				}
			}
		},
	}

	sessionRepo := repo.NewSessionRepo(s.pg)
	queryRepo := repo.NewQueryRepo(s.pg)
	responseRepo := repo.NewResponseRepo(s.pg)

	sessionController := controller.NewSessionController(sessionRepo)
	queryController := controller.NewQueryController(queryRepo)
	responseController := controller.NewResponseController(responseRepo, searchEngineConn)

	sessionWsHandler := wsHandler.NewChatHandler(sessionController, queryController, responseController)
	ws.SetupSessionSocket(s.app, sessionWsHandler, &wsConfig)

	go s.listen()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	_ = <-ch
	s.pg.Close()
	if err = searchEngineConn.Close(); err != nil {
		logger.Error(err)
	}
	if err = s.app.Shutdown(); err != nil {
		logger.Error(err)
	}
	logger.Info("shutting down the server")
}

func (s *Server) listen() {
	addr := fmt.Sprintf(":%d", s.cfg.Server.Port)
	logger.Infof("starting Server %s", addr)
	if err := s.app.Listen(addr); err != nil {
		logger.Fatalf("error has occurred while listening on %s: %v", addr, err)
	}
}
