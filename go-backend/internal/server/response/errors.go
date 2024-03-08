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
	shared.ErrWsProtocolRequired: {
		Status: http.StatusUpgradeRequired,
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
