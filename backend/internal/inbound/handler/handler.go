package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	"my-web-app.com/smart-logistic-hub/internal/inbound/dto"
	"my-web-app.com/smart-logistic-hub/internal/inbound/service"
)

type Handler struct {
	Service *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{Service: svc}
}

// @Summary Create an inbound receipt
// @Tags Inbounds
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param inbound body dto.CreateInboundRequest true "Inbound payload"
// @Success 201 {object} dto.InboundResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /inbounds [post]
func (h *Handler) Create(c *gin.Context) {
	var req dto.CreateInboundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(fmt.Errorf("%w: %v", apierrors.ErrBadRequest, err))
		return
	}

	in, err := h.Service.Create(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, dto.ToResponse(in, nil))
}

// @Summary List inbound receipts
// @Tags Inbounds
// @Produce json
// @Security BearerAuth
// @Param skip query int false "Number of items to skip" default(0)
// @Param limit query int false "Max items per page" default(20)
// @Success 200 {object} dto.PaginatedResponse
// @Failure 500 {object} map[string]interface{}
// @Router /inbounds [get]
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

// @Summary Get an inbound receipt by ID
// @Tags Inbounds
// @Produce json
// @Security BearerAuth
// @Param id path int true "Inbound ID"
// @Success 200 {object} dto.InboundResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /inbounds/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(fmt.Errorf("%w: invalid id", apierrors.ErrBadRequest))
		return
	}

	in, items, err := h.Service.Get(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToResponse(in, items))
}

// @Summary Update an inbound receipt
// @Tags Inbounds
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Inbound ID"
// @Param inbound body dto.UpdateInboundRequest true "Inbound update"
// @Success 200 {object} dto.InboundResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /inbounds/{id} [patch]
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(fmt.Errorf("%w: invalid id", apierrors.ErrBadRequest))
		return
	}

	var req dto.UpdateInboundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(fmt.Errorf("%w: %v", apierrors.ErrBadRequest, err))
		return
	}

	in, err := h.Service.Update(c.Request.Context(), id, &req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToResponse(in, nil))
}

// @Summary Delete an inbound receipt
// @Tags Inbounds
// @Security BearerAuth
// @Param id path int true "Inbound ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /inbounds/{id} [delete]
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
