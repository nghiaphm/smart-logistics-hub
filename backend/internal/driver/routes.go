package driver

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"my-web-app.com/smart-logistic-hub/internal/driver/handler"
	"my-web-app.com/smart-logistic-hub/internal/driver/repository"
	"my-web-app.com/smart-logistic-hub/internal/driver/service"
)

func RegisterRoutes(rg *gin.RouterGroup, db *sql.DB, authMw gin.HandlerFunc) {
	repo := repository.NewRepository(db)
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)
	h.RegisterRoutes(rg, authMw)
}
