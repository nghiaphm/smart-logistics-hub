package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"

	"github.com/gin-gonic/gin"
)

// RequestIDMiddleware generates or reuses a request ID, stores it in the gin
// context, attaches it to a request-scoped slog logger, and echoes it in the
// X-Request-ID response header.
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader("X-Request-ID")
		if reqID == "" {
			reqID = newRequestID()
		}

		c.Set("request_id", reqID)
		c.Header("X-Request-ID", reqID)
		c.Set("logger", slog.Default().With(slog.String("request_id", reqID)))

		c.Next()
	}
}

// LoggerFromContext returns the request-scoped slog logger set by
// RequestIDMiddleware, falling back to the default logger when absent.
func LoggerFromContext(c *gin.Context) *slog.Logger {
	if l, ok := c.Get("logger"); ok {
		if logger, ok := l.(*slog.Logger); ok {
			return logger
		}
	}
	return slog.Default()
}

func newRequestID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "unknown"
	}
	return hex.EncodeToString(b)
}
