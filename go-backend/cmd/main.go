package main

import (
	"github.com/yogenyslav/logger"
	"hack/config"
	"hack/internal/server"
)

func main() {
	cfg := config.MustNew("config/config.yaml")

	logger.SetLevel(logger.ParseLevel(cfg.Server.LogLevel))
	logger.SetFileOutput("./logs/hack-template-ssr.log", true)
	logger.Debugf("loaded config: %+v", *cfg)

	s := server.New(cfg)
	s.Run()
}
