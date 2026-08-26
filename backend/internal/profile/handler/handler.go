package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	"my-web-app.com/smart-logistic-hub/internal/profile/dto"
	"my-web-app.com/smart-logistic-hub/internal/profile/service"
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

// @Summary Get current user profile
// @Tags Profile
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.ProfileResponse
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /profile [get]
func (h *Handler) Get(c *gin.Context) {
	sub, err := userSub(c)
	if err != nil {
		c.Error(err)
		return
	}

	p, err := h.Service.Get(c.Request.Context(), sub)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToResponse(p))
}

// @Summary Create user profile
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param profile body dto.CreateProfileRequest true "Create user profile"
// @Success 201 {object} dto.ProfileResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /profile [post]
func (h *Handler) Create(c *gin.Context) {
	sub, err := userSub(c)
	if err != nil {
		c.Error(err)
		return
	}

	var req dto.CreateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(fmt.Errorf("%w: %v", apierrors.ErrBadRequest, err))
		return
	}

	p, err := h.Service.Create(c.Request.Context(), sub, &req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, dto.ToResponse(p))
}

// @Summary Update current user profile
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param profile body dto.UpdateProfileRequest true "Profile update"
// @Success 200 {object} dto.ProfileResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /profile [put]
func (h *Handler) Update(c *gin.Context) {
	sub, err := userSub(c)
	if err != nil {
		c.Error(err)
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(fmt.Errorf("%w: %v", apierrors.ErrBadRequest, err))
		return
	}

	p, err := h.Service.Update(c.Request.Context(), sub, &req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToResponse(p))
}
