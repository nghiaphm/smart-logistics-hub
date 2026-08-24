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
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"my-web-app.com/smart-logistic-hub/internal/ai"
	"my-web-app.com/smart-logistic-hub/internal/billing"
	"my-web-app.com/smart-logistic-hub/internal/driver"
	"my-web-app.com/smart-logistic-hub/internal/driver/handler"
	driverrepo "my-web-app.com/smart-logistic-hub/internal/driver/repository"
	driverservice "my-web-app.com/smart-logistic-hub/internal/driver/service"
	"my-web-app.com/smart-logistic-hub/internal/inbound"
	"my-web-app.com/smart-logistic-hub/internal/infrastructure/config"
	"my-web-app.com/smart-logistic-hub/internal/infrastructure/database"
	"my-web-app.com/smart-logistic-hub/internal/infrastructure/keycloak"
	"my-web-app.com/smart-logistic-hub/internal/infrastructure/logger"
	"my-web-app.com/smart-logistic-hub/internal/infrastructure/middleware"
	infraredis "my-web-app.com/smart-logistic-hub/internal/infrastructure/redis"
	"my-web-app.com/smart-logistic-hub/internal/inventory"
	invrepo "my-web-app.com/smart-logistic-hub/internal/inventory/repository"
	"my-web-app.com/smart-logistic-hub/internal/order"
	orderhandler "my-web-app.com/smart-logistic-hub/internal/order/handler"
	orderrepo "my-web-app.com/smart-logistic-hub/internal/order/repository"
	orderservice "my-web-app.com/smart-logistic-hub/internal/order/service"
	"my-web-app.com/smart-logistic-hub/internal/product"
	prodrepo "my-web-app.com/smart-logistic-hub/internal/product/repository"
	"my-web-app.com/smart-logistic-hub/internal/tracking"
	"my-web-app.com/smart-logistic-hub/internal/trip"
	"my-web-app.com/smart-logistic-hub/internal/warehouse"
	"my-web-app.com/smart-logistic-hub/internal/workspace"
)

// @title Smart Logistics Hub API
// @version 1.0.0
// @description Smart Logistics Hub backend API. Authentication uses the "Authorization" header with format "Bearer &lt;token&gt;" (JWT issued by Keycloak).
// @host localhost:8000
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
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
	r.Use(middleware.RequestIDMiddleware(), gin.Recovery(), corsMw, middleware.MetricsMiddleware(), middleware.ErrorHandler())

	api := r.Group("/api/v1")

	driverRepo := driverrepo.NewRepository(db)
	driverSvc := driverservice.NewService(driverRepo)
	driverHandler := handler.NewHandler(driverSvc)

	orderRepo := orderrepo.NewRepository(db)
	orderProductRepo := prodrepo.NewRepository(db)
	orderInventoryRepo := invrepo.NewRepository(db)
	orderSvc := orderservice.NewService(orderRepo, orderProductRepo, orderInventoryRepo)
	orderHandler := orderhandler.NewHandler(orderSvc)

	protected := api.Group("")
	protected.Use(authMw)
	{
		driver.RegisterRoutes(protected, driverHandler)
		order.RegisterRoutes(protected, orderHandler)
		inventory.RegisterRoutes(protected, db, authMw)
		tracking.RegisterRoutes(protected, db, authMw)
		product.RegisterRoutes(protected, db, authMw)
		warehouse.RegisterRoutes(protected, db, authMw)
		trip.RegisterRoutes(protected, db, authMw)
		inbound.RegisterRoutes(protected, db, authMw)
		billing.RegisterRoutes(protected, db, authMw)
		ai.RegisterRoutes(protected, db, authMw)
		workspace.RegisterRoutes(protected, db, authMw)
	}

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "alive"})
	})

	r.GET("/readyz", func(c *gin.Context) {
		if err := db.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready", "database": "disconnected"})
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

	// Internal metrics server (Prometheus scraping), isolated from public traffic.
	var metricsSrv *http.Server
	if cfg.MetricsEnabled {
		metricsSrv = &http.Server{
			Addr:              fmt.Sprintf("%s:%s", cfg.MetricsHost, cfg.MetricsPort),
			Handler:           promhttp.Handler(),
			ReadHeaderTimeout: 5 * time.Second,
		}
		go func() {
			log.Info("metrics server listening", "addr", metricsSrv.Addr)
			if err := metricsSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Error("metrics server failed", "error", err)
			}
		}()
	}

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
	if metricsSrv != nil {
		if err := metricsSrv.Shutdown(ctx); err != nil {
			log.Error("metrics server shutdown failed", "error", err)
		}
	}
	log.Info("server stopped")
}
