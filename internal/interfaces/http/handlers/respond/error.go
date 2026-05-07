package respond

import (
	"errors"
	"net/http"

	"github.com/Nurdiansyah15/ddd-arch/internal/app/apperror"
	"github.com/gin-gonic/gin"
)

// kindStatus memetakan Kind → HTTP status code.
// Satu-satunya tempat di seluruh codebase yang tahu HTTP status code dari error.
var kindStatus = map[apperror.Kind]int{
	apperror.KindNotFound:     http.StatusNotFound,
	apperror.KindConflict:     http.StatusConflict,
	apperror.KindUnauthorized: http.StatusUnauthorized,
	apperror.KindForbidden:    http.StatusForbidden,
	apperror.KindValidation:   http.StatusBadRequest,
	apperror.KindInternal:     http.StatusInternalServerError,
}

// Error mengubah error dari usecase menjadi JSON response yang tepat.
// Untuk KindInternal, Cause tidak pernah dikirim ke client.
func Error(c *gin.Context, err error) {
	var ae *apperror.AppError
	if !errors.As(err, &ae) {
		// error tidak terstruktur (tidak dari usecase) → internal
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	status, ok := kindStatus[ae.Kind]
	if !ok {
		status = http.StatusInternalServerError
	}

	c.JSON(status, gin.H{"error": ae.Message})
}
