package warehouse

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"my-web-app.com/smart-logistic-hub/internal/infrastructure/middleware"
	"my-web-app.com/smart-logistic-hub/internal/warehouse/handler"
	whrepo "my-web-app.com/smart-logistic-hub/internal/warehouse/repository"
	"my-web-app.com/smart-logistic-hub/internal/warehouse/service"
)

func RegisterRoutes(router *gin.RouterGroup, db *sql.DB, authMw gin.HandlerFunc) {
	repo := whrepo.NewRepository(db)
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	group := router.Group("/warehouses")
	group.Use(authMw)
	{
		group.POST("", h.Create)
		group.GET("", h.List)
		group.GET("/:id", h.Get)
		group.PATCH("/:id", h.Update)
		group.DELETE("/:id", middleware.RequireRole("admin"), h.Delete)
	}
}
