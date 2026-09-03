package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	"my-web-app.com/smart-logistic-hub/internal/workspace/dto"
	"my-web-app.com/smart-logistic-hub/internal/workspace/service"
)

type Handler struct {
	Service *service.Service
	Members *service.MembershipService
}

func NewHandler(svc *service.Service, members *service.MembershipService) *Handler {
	return &Handler{Service: svc, Members: members}
}

func currentUserSub(c *gin.Context) string {
	claims, ok := c.Get("user")
	if !ok {
		return ""
	}
	if mc, ok := claims.(jwt.MapClaims); ok {
		if sub, ok := mc["sub"].(string); ok {
			return sub
		}
	}
	return ""
}

// @Summary Create a workspace
// @Tags Workspaces
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param workspace body dto.CreateWorkspaceRequest true "Workspace payload"
// @Success 201 {object} dto.WorkspaceResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /workspaces [post]
func (h *Handler) Create(c *gin.Context) {
	var req dto.CreateWorkspaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(fmt.Errorf("%w: %v", apierrors.ErrBadRequest, err))
		return
	}

	w, err := h.Service.Create(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, dto.ToResponse(w))
}

// @Summary List workspaces
// @Tags Workspaces
// @Produce json
// @Security BearerAuth
// @Param skip query int false "Number of items to skip" default(0)
// @Param limit query int false "Max items per page" default(20)
// @Success 200 {object} dto.PaginatedResponse
// @Failure 500 {object} map[string]interface{}
// @Router /workspaces [get]
func (h *Handler) List(c *gin.Context) {
	skip, _ := strconv.Atoi(c.DefaultQuery("skip", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit <= 0 {
		limit = 20
	}
	if skip < 0 {
		skip = 0
	}

	items, total, err := h.Service.List(c.Request.Context(), skip, limit)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.PaginatedResponse{
		Items: dto.ToResponseList(items),
		Total: total,
		Skip:  skip,
		Limit: limit,
	})
}

// @Summary Get a workspace by ID
// @Tags Workspaces
// @Produce json
// @Security BearerAuth
// @Param id path int true "Workspace ID"
// @Success 200 {object} dto.WorkspaceResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /workspaces/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(fmt.Errorf("%w: invalid id", apierrors.ErrBadRequest))
		return
	}

	w, err := h.Service.Get(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToResponse(w))
}

// @Summary Update a workspace
// @Tags Workspaces
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Workspace ID"
// @Param workspace body dto.UpdateWorkspaceRequest true "Workspace update"
// @Success 200 {object} dto.WorkspaceResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /workspaces/{id} [patch]
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(fmt.Errorf("%w: invalid id", apierrors.ErrBadRequest))
		return
	}

	var req dto.UpdateWorkspaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(fmt.Errorf("%w: %v", apierrors.ErrBadRequest, err))
		return
	}

	w, err := h.Service.Update(c.Request.Context(), id, &req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToResponse(w))
}

// @Summary Delete a workspace
// @Tags Workspaces
// @Security BearerAuth
// @Param id path int true "Workspace ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /workspaces/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(fmt.Errorf("%w: invalid id", apierrors.ErrBadRequest))
		return
	}

	if err := h.Service.Delete(c.Request.Context(), id); err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

// @Summary Set / unset a user as admin of a workspace
// @Tags Workspaces
// @Description Upsert membership (workspace_id, user_id) and set is_admin.
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Workspace ID"
// @Param user_id path int true "User ID (users.id)"
// @Param body body dto.SetMemberAdminRequest true "is_admin flag"
// @Success 200 {object} dto.WorkspaceUserResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /workspaces/{id}/members/{user_id} [put]
func (h *Handler) SetMemberAdmin(c *gin.Context) {
	workspaceID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(fmt.Errorf("%w: invalid workspace id", apierrors.ErrBadRequest))
		return
	}
	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		c.Error(fmt.Errorf("%w: invalid user id", apierrors.ErrBadRequest))
		return
	}

	var req dto.SetMemberAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(fmt.Errorf("%w: %v", apierrors.ErrBadRequest, err))
		return
	}
	if req.IsAdmin == nil {
		c.Error(fmt.Errorf("%w: is_admin is required", apierrors.ErrBadRequest))
		return
	}

	wu, err := h.Members.SetIsAdmin(c.Request.Context(), workspaceID, userID, *req.IsAdmin, currentUserSub(c))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToWorkspaceUserResponse(wu))
}
