package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// BodyLimit caps the total request body at maxBytes by wrapping it in an
// http.MaxBytesReader. Reads past the limit fail with *http.MaxBytesError,
// which handlers can translate into a clear 413 via IsOversizeBody.
//
// Installing the reader is enough to stop unbounded streaming: Gin's
// c.FormFile / c.PostForm / c.ShouldBindJSON all read through c.Request.Body
// and will surface the MaxBytesError.
func BodyLimit(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if maxBytes > 0 {
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		}
		c.Next()
	}
}

// IsOversizeBody reports whether err originated from an http.MaxBytesReader
// whose limit was exceeded. Handlers use it to return 413 instead of 400/500.
func IsOversizeBody(err error) bool {
	_, ok := err.(*http.MaxBytesError)
	return ok
}
