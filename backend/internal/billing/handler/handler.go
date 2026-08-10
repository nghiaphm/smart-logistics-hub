package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"my-web-app.com/smart-logistic-hub/internal/billing/dto"
	"my-web-app.com/smart-logistic-hub/internal/billing/service"
	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
)

type Handler struct {
	Service *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{Service: svc}
}

func (h *Handler) Create(c *gin.Context) {
	var req dto.CreateBillingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(fmt.Errorf("%w: %v", apierrors.ErrBadRequest, err))
		return
	}

	b, err := h.Service.Create(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, dto.ToResponse(b))
}

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

func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(fmt.Errorf("%w: invalid id", apierrors.ErrBadRequest))
		return
	}

	b, err := h.Service.Get(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToResponse(b))
}

func (h *Handler) GetByCode(c *gin.Context) {
	b, err := h.Service.GetByCode(c.Request.Context(), c.Param("billing_code"))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToResponse(b))
}

func (h *Handler) GetByOrderCode(c *gin.Context) {
	b, err := h.Service.GetByOrderCode(c.Request.Context(), c.Param("order_code"))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToResponse(b))
}

func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(fmt.Errorf("%w: invalid id", apierrors.ErrBadRequest))
		return
	}

	var req dto.UpdateBillingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(fmt.Errorf("%w: %v", apierrors.ErrBadRequest, err))
		return
	}

	b, err := h.Service.Update(c.Request.Context(), id, &req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToResponse(b))
}

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
