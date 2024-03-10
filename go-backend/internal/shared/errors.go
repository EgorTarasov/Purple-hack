package shared

import (
	"errors"
)

// 400
var (
	ErrDuplicateKey          = errors.New("duplicate key")
	ErrWsProtocolRequired    = errors.New("upgrade to websocket protocol is required")
	ErrSessionIdInvalid      = errors.New("invalid session id value")
	ErrUnexpectedMessageType = errors.New("unexpected message type, use TextMessage instead")
)

// 500
var (
	ErrInsertRecord   = errors.New("insert record failed")
	ErrFindRecord     = errors.New("find record failed")
	ErrUpdateRecord   = errors.New("update record failed")
	ErrDeleteRecord   = errors.New("delete record failed")
	ErrCipherTooShort = errors.New("cipher text is too short")
)
