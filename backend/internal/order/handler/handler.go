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
	c.JSON(http.StatusCreated, o)
}

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

	resp := dto.OrderResponse{
		ID:                 o.ID,
		OrderCode:          o.OrderCode,
		SenderName:         o.SenderName,
		SenderPhone:        o.SenderPhone,
		SenderAddress:      o.SenderAddress,
		SenderProvince:     o.SenderProvince,
		SenderDistrict:     o.SenderDistrict,
		SenderWard:         o.SenderWard,
		SenderPostalCode:   o.SenderPostalCode,
		ReceiverName:       o.ReceiverName,
		ReceiverPhone:      o.ReceiverPhone,
		ReceiverAddress:    o.ReceiverAddress,
		ReceiverProvince:   o.ReceiverProvince,
		ReceiverDistrict:   o.ReceiverDistrict,
		ReceiverWard:       o.ReceiverWard,
		ReceiverPostalCode: o.ReceiverPostalCode,
		Status:             o.Status,
		AssignedDriverID:   o.AssignedDriverID,
		CreatedAt:          o.CreatedAt,
		UpdatedAt:          o.UpdatedAt,
		CreatedBy:          o.CreatedBy,
	}
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

func (h *Handler) List(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("skip", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 200 {
		limit = 10
	}

	orders, total, err := h.svc.List(c.Request.Context(), offset, limit)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"items": orders,
		"total": total,
		"skip":  offset,
		"limit": limit,
	})
}

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
	c.JSON(http.StatusOK, o)
}

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
