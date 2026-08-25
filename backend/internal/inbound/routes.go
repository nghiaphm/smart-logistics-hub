package inbound

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"my-web-app.com/smart-logistic-hub/internal/inbound/handler"
	inbrepo "my-web-app.com/smart-logistic-hub/internal/inbound/repository"
	"my-web-app.com/smart-logistic-hub/internal/inbound/service"
	"my-web-app.com/smart-logistic-hub/internal/infrastructure/middleware"
	invrepo "my-web-app.com/smart-logistic-hub/internal/inventory/repository"
)

func RegisterRoutes(router *gin.RouterGroup, db *sql.DB, authMw gin.HandlerFunc) {
	repo := inbrepo.NewRepository(db)
	inventoryRepo := invrepo.NewRepository(db)
	svc := service.NewService(repo, inventoryRepo)
	h := handler.NewHandler(svc)

	group := router.Group("/inbounds")
	group.Use(authMw)
	{
		group.POST("", h.Create)
		group.GET("", h.List)
		group.GET("/:id", h.Get)
		group.PATCH("/:id", h.Update)
		group.DELETE("/:id", middleware.RequireRole("system_admin"), h.Delete)
	}
}
