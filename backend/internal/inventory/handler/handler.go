package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	"my-web-app.com/smart-logistic-hub/internal/inventory/dto"
	"my-web-app.com/smart-logistic-hub/internal/inventory/service"
)

type Handler struct {
	Service *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{Service: svc}
}

// @Summary Create an inventory record
// @Tags Inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param inventory body dto.CreateInventoryRequest true "Inventory payload"
// @Success 201 {object} dto.InventoryResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /inventory [post]
func (h *Handler) Create(c *gin.Context) {
	var req dto.CreateInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(fmt.Errorf("%w: %v", apierrors.ErrBadRequest, err))
		return
	}

	inv, err := h.Service.Create(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, dto.ToResponse(inv))
}

// @Summary List inventory records
// @Tags Inventory
// @Produce json
// @Security BearerAuth
// @Param skip query int false "Number of items to skip" default(0)
// @Param limit query int false "Max items per page" default(20)
// @Success 200 {object} dto.PaginatedResponse
// @Failure 500 {object} map[string]interface{}
// @Router /inventory [get]
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

// @Summary Get an inventory record by ID
// @Tags Inventory
// @Produce json
// @Security BearerAuth
// @Param id path int true "Inventory ID"
// @Success 200 {object} dto.InventoryResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /inventory/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(fmt.Errorf("%w: invalid id", apierrors.ErrBadRequest))
		return
	}

	inv, err := h.Service.Get(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToResponse(inv))
}

// @Summary Update an inventory record
// @Tags Inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Inventory ID"
// @Param inventory body dto.UpdateInventoryRequest true "Inventory update"
// @Success 200 {object} dto.InventoryResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /inventory/{id} [patch]
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(fmt.Errorf("%w: invalid id", apierrors.ErrBadRequest))
		return
	}

	var req dto.UpdateInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(fmt.Errorf("%w: %v", apierrors.ErrBadRequest, err))
		return
	}

	inv, err := h.Service.Update(c.Request.Context(), id, &req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToResponse(inv))
}

// @Summary Delete an inventory record
// @Tags Inventory
// @Security BearerAuth
// @Param id path int true "Inventory ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /inventory/{id} [delete]
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
