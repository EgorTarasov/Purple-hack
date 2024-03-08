package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/yogenyslav/logger"
	"github.com/yogenyslav/storage/postgres"
)

type Config struct {
	Server   ServerConfig    `yaml:"server"`
	Postgres postgres.Config `yaml:"postgres"`
}

type ServerConfig struct {
	Port        int    `yaml:"port"`
	LogLevel    string `yaml:"logLevel"`
	CorsOrigins string `yaml:"corsOrigins"`
}

func MustNew(path string) *Config {
	cfg := &Config{}
	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		logger.Panic(err)
	}
	return cfg
}
