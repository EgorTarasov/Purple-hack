package data

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Id          uuid.UUID `db:"id"`
	PromptIds   []int64   `db:"prompt_ids"`
	ResponseIds []int64   `db:"response_ids"`
	CreatedAt   time.Time `db:"created_at"`
}

type Prompt struct {
	Id        int64     `db:"id"`
	SessionId uuid.UUID `db:"fk_session_id"`
	Body      string    `db:"body"`
	UserAgent string    `db:"user_agent"`
	CreatedAt time.Time `db:"created_at"`
}

type Response struct {
	Id        int64     `db:"id"`
	SessionId uuid.UUID `db:"fk_session_id"`
	PromptId  int64     `db:"fk_prompt_id"`
	Body      string    `db:"body"`
	CreatedAt time.Time `db:"created_at"`
}
