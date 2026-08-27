package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	"my-web-app.com/smart-logistic-hub/internal/role/dto"
	"my-web-app.com/smart-logistic-hub/internal/role/service"
)

type Handler struct {
	Service *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{Service: svc}
}

func userSub(c *gin.Context) (string, error) {
	raw, ok := c.Get("user")
	if !ok {
		return "", fmt.Errorf("%w: missing user claims", apierrors.ErrUnauthorized)
	}
	claims, ok := raw.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("%w: invalid user claims", apierrors.ErrUnauthorized)
	}
	sub, _ := claims["sub"].(string)
	if sub == "" {
		return "", fmt.Errorf("%w: missing sub claim", apierrors.ErrUnauthorized)
	}
	return sub, nil
}

// @Summary Get roles and permissions for current user
// @Tags Role
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.UserRolesAndPermissionsResponse
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /roles/me [get]
func (h *Handler) GetMe(c *gin.Context) {
	sub, err := userSub(c)
	if err != nil {
		c.Error(err)
		return
	}

	res, err := h.Service.GetUserRolesAndPermissions(c.Request.Context(), sub)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, res)
}
