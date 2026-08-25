package ai

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"my-web-app.com/smart-logistic-hub/internal/ai/handler"
	airepo "my-web-app.com/smart-logistic-hub/internal/ai/repository"
	"my-web-app.com/smart-logistic-hub/internal/ai/service"
	"my-web-app.com/smart-logistic-hub/internal/infrastructure/middleware"
)

func RegisterRoutes(router *gin.RouterGroup, db *sql.DB, authMw gin.HandlerFunc) {
	repo := airepo.NewRepository(db)
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	group := router.Group("/ai-events")
	group.Use(authMw)
	{
		group.POST("", h.Create)
		group.GET("", h.List)
		group.GET("/:id", h.Get)
		group.PATCH("/:id", h.Update)
		group.DELETE("/:id", middleware.RequireRole("system_admin"), h.Delete)
	}
}
