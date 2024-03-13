package main

import (
	"context"

	"purple/internal/api/data"

	"github.com/yogenyslav/logger"
	"github.com/yogenyslav/storage/postgres"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg := postgres.Config{
		User:     "pguser",
		Password: "pgpass",
		Host:     "localhost",
		Port:     5432,
		Db:       "cbr",
		Ssl:      false,
	}
	pg := postgres.MustNew(&cfg, 20)

	user := data.User{
		Email:    "misis_banach_space@test.com",
		Password: "test123456",
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Panic(err)
	}

	user.Password = string(hashedPass)
	_, err = pg.Exec(context.Background(), `
		insert into "user"(email, password)
		values ($1, $2);
	`, user.Email, user.Password)
	if err != nil {
		logger.Panic(err)
	}
}
