package server

import (
	"fmt"
	"github.com/gofiber/swagger"
	"hack/config"
	_ "hack/docs"
	"hack/internal/api/data/repo"
	"hack/internal/api/domain/controller"
	"hack/internal/api/transport/http/handler"
	"hack/internal/api/transport/http/middleware"
	"hack/internal/api/transport/http/router"
	authRepo "hack/internal/auth/data/repo"
	authController "hack/internal/auth/domain/controller"
	authHandler "hack/internal/auth/transport/http/handler"
	authRouter "hack/internal/auth/transport/http/router"
	"hack/internal/server/response"
	"hack/pkg/mailing"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	loggermw "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/yogenyslav/logger"
	"github.com/yogenyslav/storage/postgres"
)

type Server struct {
	app         *fiber.App
	cfg         *config.Config
	pg          *pgxpool.Pool
	redisClient *redis.Client
}

func New(cfg *config.Config) *Server {
	app := fiber.New(fiber.Config{
		ServerHeader: "Fiber",
		BodyLimit:    500 * 1024 * 1024,
		ErrorHandler: response.ErrorHandler,
	})

	pg := postgres.MustNew(&cfg.Postgres, 20)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: "",
		DB:       cfg.Redis.Db,
	})

	app.Use(loggermw.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Server.CorsOrigins,
		AllowMethods:     "*",
		AllowHeaders:     "*",
		AllowCredentials: true,
	}))

	app.Get("/swagger/*", swagger.HandlerDefault)

	return &Server{
		app:         app,
		cfg:         cfg,
		pg:          pg,
		redisClient: redisClient,
	}
}

func (s *Server) Run() {
	defer s.pg.Close()

	userRepo := repo.NewUserRepo(s.pg, s.redisClient)
	tokenRepo := authRepo.NewTokenRepo(s.redisClient)
	ac := authController.New(userRepo, tokenRepo, &s.cfg.Jwt)
	ah := authHandler.New(ac)
	authRouter.Setup(s.app, ah)

	mw := middleware.New(tokenRepo, &s.cfg.Jwt)
	apiRouter := s.app.Group("/api")

	mailSever := mailing.NewMailServer(&s.cfg.Mail)

	userController := controller.NewUserController(userRepo, mailSever)
	userHandler := handler.NewUserHandler(userController)
	router.SetupUserRoutes(apiRouter, userHandler, mw)

	go s.listen()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	_ = <-ch
	logger.Info("shutting down the server")
}

func (s *Server) listen() {
	addr := fmt.Sprintf(":%d", s.cfg.Server.Port)
	logger.Infof("starting Server %s", addr)
	if err := s.app.Listen(addr); err != nil {
		logger.Fatalf("error has occurred while listening on %s: %v", addr, err)
	}
}
