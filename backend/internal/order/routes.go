package order

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"my-web-app.com/smart-logistic-hub/internal/order/handler"
	"my-web-app.com/smart-logistic-hub/internal/order/service"
)

func RegisterRoutes(rg *gin.RouterGroup, db *sql.DB, authMw gin.HandlerFunc) {
	svc := service.NewService(db)
	h := handler.NewHandler(svc)
	h.RegisterRoutes(rg, authMw)
}
