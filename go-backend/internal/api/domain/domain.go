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
	Id          uuid.UUID  `json:"id"`
	QueryIds    []int64    `json:"-"`
	ResponseIds []int64    `json:"-"`
	Queries     []Query    `json:"queries"`
	Responses   []Response `json:"responses"`
	CreatedAt   time.Time  `json:"createdAt"`
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
	Context   map[string][]string
	Model     string
}

type Response struct {
	Id        int64               `json:"id"`
	Body      string              `json:"body"`
	Context   map[string][]string `json:"context"`
	CreatedAt time.Time           `json:"createdAt"`
}
