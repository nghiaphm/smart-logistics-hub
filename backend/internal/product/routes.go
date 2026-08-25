package product

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"my-web-app.com/smart-logistic-hub/internal/infrastructure/middleware"
	"my-web-app.com/smart-logistic-hub/internal/product/handler"
	prodrepo "my-web-app.com/smart-logistic-hub/internal/product/repository"
	"my-web-app.com/smart-logistic-hub/internal/product/service"
)

func RegisterRoutes(router *gin.RouterGroup, db *sql.DB, authMw gin.HandlerFunc) {
	repo := prodrepo.NewRepository(db)
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	group := router.Group("/products")
	group.Use(authMw)
	{
		group.POST("", h.Create)
		group.GET("", h.List)
		group.GET("/:id", h.Get)
		group.PATCH("/:id", h.Update)
		group.DELETE("/:id", middleware.RequireRole("system_admin"), h.Delete)
	}
}
