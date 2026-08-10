package order

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	invrepo "my-web-app.com/smart-logistic-hub/internal/inventory/repository"
	"my-web-app.com/smart-logistic-hub/internal/order/handler"
	"my-web-app.com/smart-logistic-hub/internal/order/repository"
	"my-web-app.com/smart-logistic-hub/internal/order/service"
	prodrepo "my-web-app.com/smart-logistic-hub/internal/product/repository"
)

func RegisterRoutes(rg *gin.RouterGroup, db *sql.DB, authMw gin.HandlerFunc) {
	repo := repository.NewRepository(db)
	productRepo := prodrepo.NewRepository(db)
	inventoryRepo := invrepo.NewRepository(db)
	svc := service.NewService(repo, productRepo, inventoryRepo)
	h := handler.NewHandler(svc)
	h.RegisterRoutes(rg, authMw)
}
