package config

import (
	"purple/pkg/jwt"
	"purple/pkg/mailing"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/yogenyslav/logger"
	"github.com/yogenyslav/storage/postgres"
)

type Config struct {
	Server   ServerConfig    `yaml:"server"`
	Postgres postgres.Config `yaml:"postgres"`
	Jwt      jwt.Config      `yaml:"jwt"`
	Redis    RedisConfig     `yaml:"redis"`
	Mail     mailing.Config  `yaml:"mail"`
}

type ServerConfig struct {
	Port        int    `yaml:"port"`
	LogLevel    string `yaml:"logLevel"`
	CorsOrigins string `yaml:"corsOrigins"`
}

type RedisConfig struct {
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Db       int    `yaml:"db"`
}

func MustNew(path string) *Config {
	cfg := &Config{}
	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		logger.Panic(err)
	}
	return cfg
}
