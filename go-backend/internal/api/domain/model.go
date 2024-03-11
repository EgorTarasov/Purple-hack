package domain

import (
	"time"

	"github.com/google/uuid"
)

type SessionCreate struct {
	Id    uuid.UUID
	Model string
}

type Session struct {
	Id        uuid.UUID  `json:"id"`
	Queries   []Query    `json:"queries"`
	Responses []Response `json:"responses"`
	CreatedAt time.Time  `json:"createdAt"`
}

type QueryCreate struct {
	SessionId uuid.UUID
	Model     string
	Body      string
	UserAgent string
}

type Query struct {
	Id        int64     `json:"id"`
	Model     string    `json:"model"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
}

type ResponseCreate struct {
	SessionId uuid.UUID
	QueryId   int64
	Body      string
	Context   string
	Model     string
}

type Response struct {
	Id        int64     `json:"id"`
	Body      string    `json:"body"`
	Context   string    `json:"context"`
	CreatedAt time.Time `json:"createdAt"`
}
