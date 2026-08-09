package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"my-web-app.com/smart-logistic-hub/internal/driver"
	"my-web-app.com/smart-logistic-hub/internal/infrastructure/config"
	"my-web-app.com/smart-logistic-hub/internal/infrastructure/database"
	"my-web-app.com/smart-logistic-hub/internal/infrastructure/keycloak"
	"my-web-app.com/smart-logistic-hub/internal/infrastructure/logger"
	"my-web-app.com/smart-logistic-hub/internal/infrastructure/middleware"
	infraredis "my-web-app.com/smart-logistic-hub/internal/infrastructure/redis"
	"my-web-app.com/smart-logistic-hub/internal/inventory"
	"my-web-app.com/smart-logistic-hub/internal/order"
	"my-web-app.com/smart-logistic-hub/internal/tracking"
)

func main() {
	cfg := config.LoadConfig()
	log := logger.New(cfg.Environment)
	slog.SetDefault(log)

	log.Info("starting server", "env", cfg.Environment)

	db, err := database.Connect(cfg)
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer database.Close(db)
	log.Info("database connected")

	if cfg.RedisEnabled {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		redisClient, err := infraredis.NewRedisClient(ctx, cfg)
		cancel()
		if err != nil {
			log.Warn("failed to connect to Redis", "error", err)
		} else {
			defer redisClient.Close()
			log.Info("redis connected")
		}
	}

	verifier := keycloak.NewJWTVerifier(cfg)
	authMw := middleware.AuthMiddleware(cfg, cfg.DevSkipAuth, verifier)
	corsMw := middleware.CORSMiddleware(cfg.FrontendURL)

	r := gin.New()
	r.Use(gin.Recovery(), corsMw, middleware.ErrorHandler())

	api := r.Group("/api/v1")
	{
		order.RegisterRoutes(api, db, authMw)
		driver.RegisterRoutes(api, db, authMw)
		inventory.RegisterRoutes(api, db, authMw)
		tracking.RegisterRoutes(api, db, authMw)
	}

	r.GET("/health", func(c *gin.Context) {
		if err := db.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "database": "disconnected"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	r.GET("/readiness", func(c *gin.Context) {
		if err := db.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	addr := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Info("server listening", "addr", addr)
		errCh <- srv.ListenAndServe()
	}()

	shutdownCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			log.Error("server failed", "error", err)
			os.Exit(1)
		}
	case <-shutdownCtx.Done():
		log.Info("shutdown signal received")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("graceful shutdown failed", "error", err)
	}
	log.Info("server stopped")
}
