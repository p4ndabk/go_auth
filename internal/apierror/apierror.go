// Package apierror is the shared, cross-cutting way handlers turn an error
// into an HTTP response. It has no model/service of its own, same as
// health/docs.
package apierror

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Detail is the machine/human readable payload of an error response.
type Detail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Body is the envelope every JSON error response shares.
type Body struct {
	Error Detail `json:"error"`
}

// AppError is a domain/HTTP error carrying the status code to respond with.
type AppError struct {
	Status  int
	Code    string
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func New(status int, code, message string) *AppError {
	return &AppError{Status: status, Code: code, Message: message}
}

func BadRequest(code, message string) *AppError {
	return New(http.StatusBadRequest, code, message)
}

func Unauthorized(code, message string) *AppError {
	return New(http.StatusUnauthorized, code, message)
}

func NotFound(code, message string) *AppError {
	return New(http.StatusNotFound, code, message)
}

func Conflict(code, message string) *AppError {
	return New(http.StatusConflict, code, message)
}

// Respond writes err as the shared JSON envelope. It recognizes *AppError
// and gorm.ErrRecordNotFound; anything else is logged server-side and
// answered with a generic 500 — the real error text is never sent to the
// client, to avoid leaking driver/SQL/internal detail.
func Respond(c *gin.Context, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		c.JSON(appErr.Status, Body{Error: Detail{Code: appErr.Code, Message: appErr.Message}})
		return
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, Body{Error: Detail{Code: "not_found", Message: "resource not found"}})
		return
	}

	log.Printf("unexpected error: %v", err)
	c.JSON(http.StatusInternalServerError, Body{Error: Detail{Code: "internal_error", Message: "internal server error"}})
}
