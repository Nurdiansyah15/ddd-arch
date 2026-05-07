package respond

import (
	"errors"
	"net/http"

	"github.com/Nurdiansyah15/ddd-arch/internal/app/apperror"
	"github.com/gin-gonic/gin"
)

type errorBody struct {
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	Errors     any    `json:"errors"`
}

// kindStatus memetakan Kind → HTTP status code.
var kindStatus = map[apperror.Kind]int{
	apperror.KindNotFound:     http.StatusNotFound,
	apperror.KindConflict:     http.StatusConflict,
	apperror.KindUnauthorized: http.StatusUnauthorized,
	apperror.KindForbidden:    http.StatusForbidden,
	apperror.KindValidation:   http.StatusBadRequest,
	apperror.KindInternal:     http.StatusInternalServerError,
}

// Error translates an apperror.AppError into a structured JSON error response.
// The internal cause is never sent to the client.
func Error(c *gin.Context, err error) {
	var ae *apperror.AppError
	if !errors.As(err, &ae) {
		code := http.StatusInternalServerError
		c.JSON(code, errorBody{
			Status:     "error",
			StatusCode: code,
			Errors:     "internal server error",
		})
		return
	}

	code, ok := kindStatus[ae.Kind]
	if !ok {
		code = http.StatusInternalServerError
	}

	c.JSON(code, errorBody{
		Status:     "error",
		StatusCode: code,
		Errors:     ae.Message,
	})
}
