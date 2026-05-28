package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorCode string

const (
	ErrorCodeInvalidInput  ErrorCode = "INVALID_INPUT"
	ErrorCodeInvalidID     ErrorCode = "INVALID_ID"
	ErrorCodeInvalidDate   ErrorCode = "INVALID_DATE"
	ErrorCodeUnauthorized  ErrorCode = "UNAUTHORIZED"
	ErrorCodeForbidden     ErrorCode = "FORBIDDEN"
	ErrorCodeNotFound      ErrorCode = "NOT_FOUND"
	ErrorCodeConflict      ErrorCode = "CONFLICT"
	ErrorCodeInternalError ErrorCode = "INTERNAL_ERROR"
)

func RespondError(c *gin.Context, status int, code ErrorCode, message string, details gin.H) {
	c.JSON(status, errorPayload(code, message, details))
}

func AbortError(c *gin.Context, status int, code ErrorCode, message string, details gin.H) {
	c.AbortWithStatusJSON(status, errorPayload(code, message, details))
}

func RespondInvalidRequestBody(c *gin.Context) {
	RespondInvalidInput(c, "invalid request body")
}

func RespondInvalidInput(c *gin.Context, message string) {
	RespondError(c, http.StatusBadRequest, ErrorCodeInvalidInput, message, nil)
}

func RespondInvalidID(c *gin.Context, message string) {
	RespondError(c, http.StatusBadRequest, ErrorCodeInvalidID, message, nil)
}

func RespondInvalidDate(c *gin.Context, message string) {
	RespondError(c, http.StatusBadRequest, ErrorCodeInvalidDate, message, nil)
}

func RespondUnauthorized(c *gin.Context, message string) {
	RespondError(c, http.StatusUnauthorized, ErrorCodeUnauthorized, message, nil)
}

func RespondForbidden(c *gin.Context, message string) {
	RespondError(c, http.StatusForbidden, ErrorCodeForbidden, message, nil)
}

func RespondNotFound(c *gin.Context, message string) {
	RespondError(c, http.StatusNotFound, ErrorCodeNotFound, message, nil)
}

func RespondConflict(c *gin.Context, message string) {
	RespondError(c, http.StatusConflict, ErrorCodeConflict, message, nil)
}

func RespondInternal(c *gin.Context) {
	RespondError(c, http.StatusInternalServerError, ErrorCodeInternalError, "internal server error", nil)
}

func AbortUnauthorized(c *gin.Context, message string) {
	AbortError(c, http.StatusUnauthorized, ErrorCodeUnauthorized, message, nil)
}

func AbortForbidden(c *gin.Context, message string) {
	AbortError(c, http.StatusForbidden, ErrorCodeForbidden, message, nil)
}

func AbortInternal(c *gin.Context) {
	AbortError(c, http.StatusInternalServerError, ErrorCodeInternalError, "internal server error", nil)
}

func errorPayload(code ErrorCode, message string, details gin.H) gin.H {
	if details == nil {
		details = gin.H{}
	}
	return gin.H{
		"error": gin.H{
			"code":    code,
			"message": message,
			"details": details,
		},
	}
}
