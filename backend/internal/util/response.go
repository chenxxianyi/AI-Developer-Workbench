package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response is the unified API response structure.
type Response struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Data      any    `json:"data,omitempty"`
	Error     string `json:"error,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

// Error codes matching the frontend contract.
const (
	CodeSuccess         = 0
	CodeInvalidRequest  = 40001
	CodeUploadFailed    = 40002
	CodeUnsupportedType = 40003
	CodeInvalidToolType = 40004
	CodeFileTooLarge    = 40005
	CodeUnsafeArchive   = 40006
	CodeReportNotFound  = 40401
	CodeStateConflict   = 40901
	CodeInternalError   = 50001
	CodeAIProviderError = 50002
	CodeDatabaseError   = 50003
)

// Error messages for each error code.
var errorMessages = map[int]string{
	CodeSuccess:         "success",
	CodeInvalidRequest:  "invalid request",
	CodeUploadFailed:    "upload failed",
	CodeUnsupportedType: "unsupported file type",
	CodeInvalidToolType: "invalid tool type",
	CodeFileTooLarge:    "file too large",
	CodeUnsafeArchive:   "unsafe archive",
	CodeReportNotFound:  "report not found",
	CodeStateConflict:   "report state conflict",
	CodeInternalError:   "internal server error",
	CodeAIProviderError: "ai provider error",
	CodeDatabaseError:   "database error",
}

// SuccessResponse returns a success response with data.
func SuccessResponse(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: errorMessages[CodeSuccess],
		Data:    data,
	})
}

// ErrorResponse returns an error response.
func ErrorResponse(c *gin.Context, httpStatus int, bizCode int, errDetail string) {
	requestID := c.GetString("request_id")
	c.JSON(httpStatus, Response{
		Code:      bizCode,
		Message:   errorMessages[bizCode],
		Error:     errDetail,
		RequestID: requestID,
	})
}

// BadRequest returns a 400 error with invalid request code.
func BadRequest(c *gin.Context, errDetail string) {
	ErrorResponse(c, http.StatusBadRequest, CodeInvalidRequest, errDetail)
}

// NotFound returns a 404 error.
func NotFound(c *gin.Context, errDetail string) {
	ErrorResponse(c, http.StatusNotFound, CodeReportNotFound, errDetail)
}

// InternalError returns a 500 error.
func InternalError(c *gin.Context, errDetail string) {
	ErrorResponse(c, http.StatusInternalServerError, CodeInternalError, errDetail)
}

// WriteDownloadResponse writes a file download response with proper headers.
func WriteDownloadResponse(c *gin.Context, filename string, content []byte, mimeType string) {
	c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
	c.Header("Content-Type", mimeType)
	c.Data(http.StatusOK, mimeType, content)
}