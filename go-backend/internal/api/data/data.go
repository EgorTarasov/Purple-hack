package data

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Id          uuid.UUID `db:"id"`
	QueryIds    []int64   `db:"query_ids"`
	ResponseIds []int64   `db:"response_ids"`
	CreatedAt   time.Time `db:"created_at"`
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
	Id        int64               `db:"id"`
	SessionId uuid.UUID           `db:"fk_session_id"`
	QueryId   int64               `db:"fk_query_id"`
	Body      string              `db:"body"`
	Context   map[string][]string `db:"context"`
	CreatedAt time.Time           `db:"created_at"`
}
