package shared

import (
	"errors"
)

var ErrWsProtocolRequired = errors.New("upgrade to websocket protocol is required")

// 400
var (
	ErrDuplicateKey          = errors.New("duplicate key")
	ErrSessionIdInvalid      = errors.New("invalid session id value")
	ErrUnexpectedMessageType = errors.New("unexpected message type, use TextMessage instead")
	ErrUserIdInvalid         = errors.New("invalid user id value")
)

// 401
var (
	ErrNoSuchUser      = errors.New("no user with requested email found")
	ErrInvalidPassword = errors.New("invalid password")
)

// 500
var (
	ErrInsertRecord   = errors.New("insert record failed")
	ErrFindRecord     = errors.New("find record failed")
	ErrUpdateRecord   = errors.New("update record failed")
	ErrDeleteRecord   = errors.New("delete record failed")
	ErrCipherTooShort = errors.New("cipher text is too short")
)
