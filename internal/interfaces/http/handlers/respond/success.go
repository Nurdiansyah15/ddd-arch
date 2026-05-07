package respond

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type successBody struct {
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Data       any    `json:"data"`
	Meta       any    `json:"meta,omitempty"`
}

// Success sends a standardised success response.
// Pass meta as nil when pagination / extra metadata is not needed.
func Success(c *gin.Context, code int, message string, data any, meta any) {
	c.JSON(code, successBody{
		Status:     "success",
		StatusCode: code,
		Message:    message,
		Data:       data,
		Meta:       meta,
	})
}

// OK is a convenience wrapper for 200 responses.
func OK(c *gin.Context, message string, data any) {
	Success(c, http.StatusOK, message, data, nil)
}

// Created is a convenience wrapper for 201 responses.
func Created(c *gin.Context, message string, data any) {
	Success(c, http.StatusCreated, message, data, nil)
}
