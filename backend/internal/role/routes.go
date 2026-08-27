package role

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"my-web-app.com/smart-logistic-hub/internal/role/handler"
	"my-web-app.com/smart-logistic-hub/internal/role/repository"
	"my-web-app.com/smart-logistic-hub/internal/role/service"
)

func RegisterRoutes(router *gin.RouterGroup, db *sql.DB, authMw gin.HandlerFunc) {
	repo := repository.NewRepository(db)
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	group := router.Group("/roles")
	group.Use(authMw)
	{
		group.GET("/me", h.GetMe)
	}
}
