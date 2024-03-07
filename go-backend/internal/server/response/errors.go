package response

import (
	"net/http"
	"purple/internal/shared"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

var errStatus = map[error]ErrorResponse{
	pgx.ErrNoRows: {
		Msg:    "no rows were found",
		Status: http.StatusNotFound,
	},
	fiber.ErrUnprocessableEntity: {
		Msg:    "data validation error",
		Status: http.StatusUnprocessableEntity,
	},
	// 400
	shared.ErrDuplicateKey: {
		Status: http.StatusBadRequest,
	},
	shared.ErrPasswordTooLong: {
		Status: http.StatusBadRequest,
	},
	shared.ErrPassRecentlyReset: {
		Status: http.StatusBadRequest,
	},
	shared.ErrInvalidTOTP: {
		Status: http.StatusBadRequest,
	},
	// 401
	shared.ErrNoSuchUser: {
		Status: http.StatusUnauthorized,
	},
	shared.ErrInvalidPassword: {
		Status: http.StatusUnauthorized,
	},
	shared.ErrJwtMalformed: {
		Status: http.StatusUnauthorized,
	},
	shared.ErrJwtInvalid: {
		Status: http.StatusUnauthorized,
	},
	shared.ErrJwtMissing: {
		Status: http.StatusUnauthorized,
	},
	shared.ErrJwtCreate: {
		Status: http.StatusUnauthorized,
	},
	// 403
	shared.ErrNoTOTP: {
		Status: http.StatusForbidden,
	},
	shared.ErrResetPassTimeout: {
		Status: http.StatusForbidden,
	},
	// 500
	shared.ErrInsertRecord: {
		Status: http.StatusInternalServerError,
	},
	shared.ErrFindRecord: {
		Status: http.StatusInternalServerError,
	},
	shared.ErrUpdateRecord: {
		Status: http.StatusInternalServerError,
	},
	shared.ErrDeleteRecord: {
		Status: http.StatusInternalServerError,
	},
	shared.ErrCipherTooShort: {
		Status: http.StatusInternalServerError,
	},
}
