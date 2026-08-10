package trip

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"my-web-app.com/smart-logistic-hub/internal/driver/repository"
	"my-web-app.com/smart-logistic-hub/internal/infrastructure/middleware"
	"my-web-app.com/smart-logistic-hub/internal/trip/handler"
	triprepo "my-web-app.com/smart-logistic-hub/internal/trip/repository"
	"my-web-app.com/smart-logistic-hub/internal/trip/service"
)

func RegisterRoutes(router *gin.RouterGroup, db *sql.DB, authMw gin.HandlerFunc) {
	repo := triprepo.NewRepository(db)
	driverRepo := repository.NewRepository(db)
	svc := service.NewService(repo, driverRepo)
	h := handler.NewHandler(svc)

	group := router.Group("/trips")
	group.Use(authMw)
	{
		group.POST("", h.Create)
		group.GET("", h.List)
		group.GET("/:id", h.Get)
		group.PATCH("/:id", h.Update)
		group.POST("/:id/assign-driver", h.AssignDriver)
		group.DELETE("/:id", middleware.RequireRole("admin"), h.Delete)
	}
}
