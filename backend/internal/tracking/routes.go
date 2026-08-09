package tracking

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"my-web-app.com/smart-logistic-hub/internal/infrastructure/middleware"
)

func RegisterRoutes(router *gin.RouterGroup, db *sql.DB, authMw gin.HandlerFunc) {
	repo := NewRepository(db)
	svc := NewService(repo)
	handler := NewHandler(svc)

	group := router.Group("/tracking-logs")
	group.Use(authMw)
	{
		group.POST("", handler.Create)
		group.GET("", handler.List)
		group.GET("/order/:order_code", handler.GetByOrder)
		group.GET("/:id", handler.Get)
		group.PUT("/:id", handler.Update)
		group.DELETE("/:id", middleware.RequireRole("admin"), handler.Delete)
	}
}
