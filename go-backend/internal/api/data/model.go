package data

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Id        uuid.UUID `db:"id"`
	CreatedAt time.Time `db:"created_at"`
}

type Query struct {
	Id        int64     `db:"id"`
	SessionId uuid.UUID `db:"fk_session_id"`
	Body      string    `db:"body"`
	Model     string    `db:"model"`
	UserAgent string    `db:"user_agent"`
	CreatedAt time.Time `db:"created_at"`
}

type Response struct {
	Id        int64     `db:"id"`
	SessionId uuid.UUID `db:"fk_session_id"`
	QueryId   int64     `db:"fk_query_id"`
	Body      string    `db:"body"`
	Context   string    `db:"context"`
	CreatedAt time.Time `db:"created_at"`
}

type User struct {
	Id        int64     `db:"id"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
}
