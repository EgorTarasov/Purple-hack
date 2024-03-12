package response

import (
	"errors"
	"fmt"
	"net/http"

	"purple/pkg"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/yogenyslav/logger"
)

type ErrorResponse struct {
	Msg    string `json:"msg"`
	Status int    `json:"-"`
}

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	var (
		ok               bool
		e                ErrorResponse
		validationErrors validator.ValidationErrors
	)

	if pkg.CheckErrPageNotFound(err) {
		e = ErrorResponse{
			Msg:    fmt.Sprintf("page \"%s\" not found", ctx.Path()),
			Status: http.StatusNotFound,
		}
	} else {
		logger.Error(err)
		if errors.As(err, &validationErrors) {
			e = ErrorResponse{
				Msg:    err.Error(),
				Status: http.StatusUnprocessableEntity,
			}
		} else {
			e, ok = errStatus[err]
			if !ok {
				e = ErrorResponse{
					Msg:    err.Error(),
					Status: http.StatusInternalServerError,
				}
			}
			if e.Msg == "" {
				e.Msg = err.Error()
			}
		}
	}
	return ctx.Status(e.Status).JSON(e)
}
