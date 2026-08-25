package inventory

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"my-web-app.com/smart-logistic-hub/internal/infrastructure/middleware"
	"my-web-app.com/smart-logistic-hub/internal/inventory/handler"
	invrepo "my-web-app.com/smart-logistic-hub/internal/inventory/repository"
	"my-web-app.com/smart-logistic-hub/internal/inventory/service"
)

func RegisterRoutes(router *gin.RouterGroup, db *sql.DB, authMw gin.HandlerFunc) {
	repo := invrepo.NewRepository(db)
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	group := router.Group("/inventory")
	group.Use(authMw)
	{
		group.POST("", h.Create)
		group.GET("", h.List)
		group.GET("/:id", h.Get)
		group.PATCH("/:id", h.Update)
		group.DELETE("/:id", middleware.RequireRole("system_admin"), h.Delete)
	}
}
