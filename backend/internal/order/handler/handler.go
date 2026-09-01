package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	"my-web-app.com/smart-logistic-hub/internal/order/dto"
	"my-web-app.com/smart-logistic-hub/internal/order/service"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

// @Summary Create an order
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param order body dto.CreateOrderRequest true "Order payload"
// @Success 201 {object} dto.OrderResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /orders [post]
func (h *Handler) Create(c *gin.Context) {
	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(fmt.Errorf("%w: %v", apierrors.ErrBadRequest, err))
		return
	}
	o, err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, dto.ToResponse(o))
}

// @Summary Get an order by ID
// @Tags Orders
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Success 200 {object} dto.OrderResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /orders/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(fmt.Errorf("%w: invalid order ID", apierrors.ErrBadRequest))
		return
	}
	o, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	items, err := h.svc.GetItems(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	resp := dto.ToResponse(o)
	for _, item := range items {
		resp.Items = append(resp.Items, dto.OrderItemResponse{
			ID:          item.ID,
			OrderID:     item.OrderID,
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			WeightGram:  item.WeightGram,
		})
	}
	c.JSON(http.StatusOK, resp)
}

// @Summary List orders
// @Tags Orders
// @Produce json
// @Security BearerAuth
// @Param skip query int false "Number of items to skip" default(0)
// @Param limit query int false "Max items per page" default(10)
// @Param workspace_id query int false "Filter by sender workspace ID (chỉ trả đơn của workspace đó)"
// @Success 200 {object} dto.PaginatedResponse
// @Failure 500 {object} map[string]interface{}
// @Router /orders [get]
func (h *Handler) List(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("skip", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 200 {
		limit = 10
	}

	var workspaceID *int64
	if raw := c.Query("workspace_id"); raw != "" {
		if id, err := strconv.ParseInt(raw, 10, 64); err == nil && id > 0 {
			workspaceID = &id
		}
	}

	orders, total, err := h.svc.List(c.Request.Context(), offset, limit, workspaceID)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.PaginatedResponse{
		Items: dto.ToResponseList(orders),
		Total: total,
		Skip:  offset,
		Limit: limit,
	})
}

// @Summary Update an order
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Param order body dto.UpdateOrderRequest true "Order update"
// @Success 200 {object} dto.OrderResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /orders/{id} [patch]
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(fmt.Errorf("%w: invalid order ID", apierrors.ErrBadRequest))
		return
	}
	var req dto.UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(fmt.Errorf("%w: %v", apierrors.ErrBadRequest, err))
		return
	}
	o, err := h.svc.Update(c.Request.Context(), id, &req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToResponse(o))
}

// @Summary Delete an order
// @Tags Orders
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /orders/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(fmt.Errorf("%w: invalid order ID", apierrors.ErrBadRequest))
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}
