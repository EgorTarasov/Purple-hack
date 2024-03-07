package shared

import "time"

const (
	UserIdField    = "userId"
	ExpTokenField  = "exp"
	RequestIdField = "requestId"
	TimestampField = "timestamp"

	TokenTypeBearer   = "Bearer"
	AccessTokenCookie = "accessToken"

	JwtTokenExp          = time.Hour
	PasswordResetExp     = time.Minute * 5
	PasswordResetTimeout = time.Minute
)
