package billing

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"my-web-app.com/smart-logistic-hub/internal/billing/handler"
	billrepo "my-web-app.com/smart-logistic-hub/internal/billing/repository"
	"my-web-app.com/smart-logistic-hub/internal/billing/service"
	"my-web-app.com/smart-logistic-hub/internal/infrastructure/middleware"
	"my-web-app.com/smart-logistic-hub/internal/order/repository"
)

func RegisterRoutes(router *gin.RouterGroup, db *sql.DB, authMw gin.HandlerFunc) {
	repo := billrepo.NewRepository(db)
	orderRepo := repository.NewRepository(db)
	svc := service.NewService(repo, orderRepo)
	h := handler.NewHandler(svc)

	group := router.Group("/billing")
	group.Use(authMw)
	{
		group.POST("", h.Create)
		group.GET("", h.List)
		group.GET("/code/:billing_code", h.GetByCode)
		group.GET("/order/:order_code", h.GetByOrderCode)
		group.GET("/:id", h.Get)
		group.PATCH("/:id", h.Update)
		group.DELETE("/:id", middleware.RequireRole("system_admin"), h.Delete)
	}
}
