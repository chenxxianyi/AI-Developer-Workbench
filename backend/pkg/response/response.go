package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse is the unified API response envelope.
type APIResponse struct {
	Data  interface{} `json:"data,omitempty"`
	Meta  *Meta       `json:"meta,omitempty"`
	Error *APIError   `json:"error,omitempty"`
}

// Meta holds pagination info.
type Meta struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
}

// APIError holds error details.
type APIError struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Details []ErrorDetail `json:"details,omitempty"`
}

// ErrorDetail holds per-field validation errors.
type ErrorDetail struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

// Standard error codes (per A-03 API contract).
const (
	CodeValidation       = "VALIDATION_ERROR"
	CodeUnauthorized     = "UNAUTHORIZED"
	CodeForbidden        = "FORBIDDEN"
	CodeNotFound         = "NOT_FOUND"
	CodeConflict         = "CONFLICT"
	CodePayloadTooLarge  = "PAYLOAD_TOO_LARGE"
	CodeBusinessError    = "BUSINESS_ERROR"
	CodeRateLimited      = "RATE_LIMITED"
	CodeInternalError    = "INTERNAL_ERROR"
	CodeServiceUnavail   = "SERVICE_UNAVAILABLE"
)

// Success responds with 200 + data.
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{Data: data})
}

// SuccessWithMeta responds with 200 + data + pagination meta.
func SuccessWithMeta(c *gin.Context, data interface{}, meta Meta) {
	c.JSON(http.StatusOK, APIResponse{Data: data, Meta: &meta})
}

// Created responds with 201 + data.
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, APIResponse{Data: data})
}

// Error responds with an error code and message.
func Error(c *gin.Context, httpStatus int, code, message string) {
	c.AbortWithStatusJSON(httpStatus, APIResponse{
		Error: &APIError{Code: code, Message: message},
	})
}

// ErrorWithDetails responds with an error + per-field details.
func ErrorWithDetails(c *gin.Context, httpStatus int, code, message string, details []ErrorDetail) {
	c.AbortWithStatusJSON(httpStatus, APIResponse{
		Error: &APIError{Code: code, Message: message, Details: details},
	})
}

// ValidationError is a shorthand for 400 validation errors.
func ValidationError(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, CodeValidation, message)
}

// Unauthorized is a shorthand for 401.
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, CodeUnauthorized, message)
}

// Forbidden is a shorthand for 403.
func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, CodeForbidden, message)
}

// NotFound is a shorthand for 404.
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, CodeNotFound, message)
}

// Conflict is a shorthand for 409.
func Conflict(c *gin.Context, message string) {
	Error(c, http.StatusConflict, CodeConflict, message)
}

// BusinessError is a shorthand for 422.
func BusinessError(c *gin.Context, message string) {
	Error(c, http.StatusUnprocessableEntity, CodeBusinessError, message)
}

// InternalError is a shorthand for 500.
func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, CodeInternalError, message)
}
