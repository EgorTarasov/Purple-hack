package shared

import "errors"

// 400
var (
	ErrDuplicateKey      = errors.New("duplicate key")
	ErrPasswordTooLong   = errors.New("password max length exceeded")
	ErrPassRecentlyReset = errors.New("password has been already reset recently")
	ErrInvalidTOTP       = errors.New("invalid totp passcode")
)

// 401
var (
	ErrNoSuchUser      = errors.New("no user with such email")
	ErrInvalidPassword = errors.New("invalid password")
	ErrJwtMalformed    = errors.New("jwt token is malformed")
	ErrJwtInvalid      = errors.New("invalid jwt token")
	ErrJwtMissing      = errors.New("jwt token is missing")
	ErrJwtCreate       = errors.New("failed to create jwt token")
)

// 403
var (
	ErrNoTOTP           = errors.New("access denied, no required totp found")
	ErrResetPassTimeout = errors.New("totp is not verified or time to reset password has expired")
)

// 500
var (
	ErrInsertRecord   = errors.New("insert record failed")
	ErrFindRecord     = errors.New("find record failed")
	ErrUpdateRecord   = errors.New("update record failed")
	ErrDeleteRecord   = errors.New("delete record failed")
	ErrCipherTooShort = errors.New("cipher text is too short")
)
