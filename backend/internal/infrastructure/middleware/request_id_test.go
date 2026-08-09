package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRequestIDMiddlewareGeneratesID(t *testing.T) {
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	reqID := rec.Header().Get("X-Request-ID")
	if reqID == "" {
		t.Error("X-Request-ID response header should be set")
	}
	if len(reqID) != 16 {
		t.Errorf("generated request ID length = %d, want 16 (8 bytes hex)", len(reqID))
	}
}

func TestRequestIDMiddlewareReusesIncomingHeader(t *testing.T) {
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Request-ID", "client-supplied-id-123")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if got := rec.Header().Get("X-Request-ID"); got != "client-supplied-id-123" {
		t.Errorf("X-Request-ID = %q, want %q", got, "client-supplied-id-123")
	}
}

func TestRequestIDStoredInGinContext(t *testing.T) {
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/test", func(c *gin.Context) {
		reqID, ok := c.Get("request_id")
		if !ok {
			t.Error("request_id should be stored in gin context")
			return
		}
		if reqID == "" {
			t.Error("request_id should be non-empty")
		}
		if reqID != c.Writer.Header().Get("X-Request-ID") {
			t.Errorf("gin context request_id = %v, header = %q, want equal", reqID, c.Writer.Header().Get("X-Request-ID"))
		}
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
}

func TestLoggerFromContext(t *testing.T) {
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/test", func(c *gin.Context) {
		l := LoggerFromContext(c)
		if l == nil {
			t.Error("LoggerFromContext should return a logger")
		}
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
}

func TestLoggerFromContextFallsBackToDefault(t *testing.T) {
	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		l := LoggerFromContext(c)
		if l == nil {
			t.Error("LoggerFromContext without middleware should return default logger")
		}
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
}
