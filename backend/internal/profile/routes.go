package profile

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"my-web-app.com/smart-logistic-hub/internal/profile/handler"
	prepo "my-web-app.com/smart-logistic-hub/internal/profile/repository"
	"my-web-app.com/smart-logistic-hub/internal/profile/service"
)

func RegisterRoutes(router *gin.RouterGroup, db *sql.DB, authMw gin.HandlerFunc) {
	repo := prepo.NewRepository(db)
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	group := router.Group("/profile")
	group.Use(authMw)
	{
		group.GET("", h.Get)
		group.PUT("", h.Update)
	}
}
