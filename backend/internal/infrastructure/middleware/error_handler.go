package middleware

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
)

type errorResponse struct {
	Error errorBody `json:"error"`
}

type errorBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ErrorHandler is a global middleware that renders errors recorded via c.Error.
// It must be registered with router.Use(...) so that it runs after handlers.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		var apiErr *apierrors.APIError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, errorResponse{
				Error: errorBody{Code: apiErr.StatusCode, Message: apiErr.Message},
			})
			return
		}

		// Unknown/internal errors: never leak the underlying message to clients.
		slog.Error("unhandled error",
			"path", c.Request.URL.Path,
			"method", c.Request.Method,
			"error", err.Error(),
		)
		c.JSON(http.StatusInternalServerError, errorResponse{
			Error: errorBody{Code: http.StatusInternalServerError, Message: "Internal server error"},
		})
	}
}
