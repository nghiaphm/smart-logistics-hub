package workspace

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"my-web-app.com/smart-logistic-hub/internal/infrastructure/middleware"
	userrepo "my-web-app.com/smart-logistic-hub/internal/user/repository"
	"my-web-app.com/smart-logistic-hub/internal/workspace/handler"
	wsrepo "my-web-app.com/smart-logistic-hub/internal/workspace/repository"
	"my-web-app.com/smart-logistic-hub/internal/workspace/service"
)

func RegisterRoutes(router *gin.RouterGroup, db *sql.DB, authMw gin.HandlerFunc) {
	repo := wsrepo.NewRepository(db)
	members := wsrepo.NewMemberRepository(db)
	users := userrepo.NewRepository(db)
	svc := service.NewService(repo)
	memberSvc := service.NewMembershipService(members, repo, users)
	h := handler.NewHandler(svc, memberSvc)

	group := router.Group("/workspaces")
	group.Use(authMw)
	{
		group.POST("", h.Create)
		group.GET("", h.List)
		group.GET("/:id", h.Get)
		group.PATCH("/:id", h.Update)
		group.DELETE("/:id", middleware.RequireRole("system_admin"), h.Delete)
		// Gán/thu hồi is_admin — CHỈ Super Admin (system_admin) (TASK-097).
		group.PUT("/:id/members/:user_id", middleware.RequireRole("system_admin"), h.SetMemberAdmin)
	}
}
