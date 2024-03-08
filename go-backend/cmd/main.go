package main

import (
	"purple/config"
	"purple/internal/server"

	"github.com/yogenyslav/logger"
)

func main() {
	cfg := config.MustNew("config/config.yaml")

	logger.SetLevel(logger.ParseLevel(cfg.Server.LogLevel))
	logger.SetFileOutput("./logs/purple.log", true)
	logger.Debugf("loaded config: %+v", *cfg)

	s := server.New(cfg)
	s.Run()
}
