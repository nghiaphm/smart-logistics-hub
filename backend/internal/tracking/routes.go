package tracking

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"my-web-app.com/smart-logistic-hub/internal/infrastructure/middleware"
	"my-web-app.com/smart-logistic-hub/internal/tracking/handler"
	trkrepo "my-web-app.com/smart-logistic-hub/internal/tracking/repository"
	"my-web-app.com/smart-logistic-hub/internal/tracking/service"
)

func RegisterRoutes(router *gin.RouterGroup, db *sql.DB, authMw gin.HandlerFunc) {
	repo := trkrepo.NewRepository(db)
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	group := router.Group("/tracking-logs")
	group.Use(authMw)
	{
		group.POST("", h.Create)
		group.GET("", h.List)
		group.GET("/order/:order_code", h.GetByOrder)
		group.GET("/:id", h.Get)
		group.PUT("/:id", h.Update)
		group.DELETE("/:id", middleware.RequireRole("admin"), h.Delete)
	}
}
