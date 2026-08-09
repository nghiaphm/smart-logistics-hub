package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestErrorHandlerRendersAPIError(t *testing.T) {
	router := gin.New()
	router.Use(ErrorHandler())
	router.GET("/test", func(c *gin.Context) {
		c.Error(apierrors.ErrNotFound)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}

	var body struct {
		Error struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Error.Code != http.StatusNotFound {
		t.Errorf("error.code = %d, want %d", body.Error.Code, http.StatusNotFound)
	}
	if body.Error.Message != "Not Found" {
		t.Errorf("error.message = %q, want %q", body.Error.Message, "Not Found")
	}
}

func TestErrorHandlerRendersWrappedAPIError(t *testing.T) {
	router := gin.New()
	router.Use(ErrorHandler())
	router.GET("/test", func(c *gin.Context) {
		c.Error(fmt.Errorf("%w: invalid id", apierrors.ErrBadRequest))
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestErrorHandlerMasksInternalError(t *testing.T) {
	router := gin.New()
	router.Use(ErrorHandler())
	router.GET("/test", func(c *gin.Context) {
		c.Error(fmt.Errorf("connection refused: secret dsn detail"))
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}

	body := rec.Body.String()
	if strings.Contains(body, "connection refused") || strings.Contains(body, "secret dsn") {
		t.Errorf("response leaked internal error details: %s", body)
	}

	var parsed struct {
		Error struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &parsed); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if parsed.Error.Message != "Internal server error" {
		t.Errorf("error.message = %q, want %q", parsed.Error.Message, "Internal server error")
	}
}

func TestErrorHandlerDoesNothingWithoutErrors(t *testing.T) {
	router := gin.New()
	router.Use(ErrorHandler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}
