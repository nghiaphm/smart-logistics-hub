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

// @Summary Create a billing record
// @Tags Billing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param billing body dto.CreateBillingRequest true "Billing payload"
// @Success 201 {object} dto.BillingResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /billing [post]
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

// @Summary List billing records
// @Tags Billing
// @Produce json
// @Security BearerAuth
// @Param skip query int false "Number of items to skip" default(0)
// @Param limit query int false "Max items per page" default(20)
// @Success 200 {object} dto.PaginatedResponse
// @Failure 500 {object} map[string]interface{}
// @Router /billing [get]
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

// @Summary Get a billing record by ID
// @Tags Billing
// @Produce json
// @Security BearerAuth
// @Param id path int true "Billing ID"
// @Success 200 {object} dto.BillingResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /billing/{id} [get]
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

// @Summary Get a billing record by billing code
// @Tags Billing
// @Produce json
// @Security BearerAuth
// @Param billing_code path string true "Billing code"
// @Success 200 {object} dto.BillingResponse
// @Failure 404 {object} map[string]interface{}
// @Router /billing/code/{billing_code} [get]
func (h *Handler) GetByCode(c *gin.Context) {
	b, err := h.Service.GetByCode(c.Request.Context(), c.Param("billing_code"))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToResponse(b))
}

// @Summary Get a billing record by order code
// @Tags Billing
// @Produce json
// @Security BearerAuth
// @Param order_code path string true "Order code"
// @Success 200 {object} dto.BillingResponse
// @Failure 404 {object} map[string]interface{}
// @Router /billing/order/{order_code} [get]
func (h *Handler) GetByOrderCode(c *gin.Context) {
	b, err := h.Service.GetByOrderCode(c.Request.Context(), c.Param("order_code"))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToResponse(b))
}

// @Summary Update a billing record
// @Tags Billing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Billing ID"
// @Param billing body dto.UpdateBillingRequest true "Billing update"
// @Success 200 {object} dto.BillingResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /billing/{id} [patch]
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

// @Summary Delete a billing record
// @Tags Billing
// @Security BearerAuth
// @Param id path int true "Billing ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /billing/{id} [delete]
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
