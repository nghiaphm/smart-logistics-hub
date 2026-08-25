package driver

import (
	"github.com/gin-gonic/gin"

	"my-web-app.com/smart-logistic-hub/internal/driver/handler"
	"my-web-app.com/smart-logistic-hub/internal/infrastructure/middleware"
)

func RegisterRoutes(rg *gin.RouterGroup, h *handler.Handler) {
	group := rg.Group("/drivers")
	{
		group.POST("", h.Create)
		group.GET("", h.List)
		group.GET("/:id", h.Get)
		group.PATCH("/:id", h.Update)
		group.DELETE("/:id", middleware.RequireRole("system_admin"), h.Delete)
	}
}
