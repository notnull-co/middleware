package handler

import (
	"net/http"

	"github.com/labstack/echo"
)

type CustomError interface {
	Code() int
	Message() string
	Details() interface{}
}

type HttpError struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Details interface{} `json:"details,omitempty"`
}

func Error(codeToStatus map[int]int) func(err error, ctx echo.Context) {
	return func(err error, ctx echo.Context) {
		if ctx.Response().Committed {
			return
		}
		status := http.StatusInternalServerError

		if customError, ok := err.(CustomError); ok {
			if codeToStatus != nil {
				if st, ok := codeToStatus[customError.Code()]; ok {
					status = st
				}
			}
			ctx.JSON(status, HttpError{
				Message: customError.Message(),
				Code:    customError.Code(),
				Details: customError.Details(),
			})
		} else {
			ctx.NoContent(status)
		}

		ctx.Response().Committed = true
	}
}
