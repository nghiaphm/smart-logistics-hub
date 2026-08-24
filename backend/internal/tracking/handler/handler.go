package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	"my-web-app.com/smart-logistic-hub/internal/tracking/dto"
	"my-web-app.com/smart-logistic-hub/internal/tracking/service"
)

type Handler struct {
	Service *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{Service: svc}
}

// @Summary Create a tracking event
// @Tags Tracking
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param event body dto.CreateTrackingEventRequest true "Tracking event payload"
// @Success 201 {object} dto.TrackingEventResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /tracking-logs [post]
func (h *Handler) Create(c *gin.Context) {
	var req dto.CreateTrackingEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(fmt.Errorf("%w: %v", apierrors.ErrBadRequest, err))
		return
	}
	event, err := h.Service.Create(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, dto.ToResponse(event))
}

// @Summary List tracking events
// @Tags Tracking
// @Produce json
// @Security BearerAuth
// @Param order_code query string false "Filter by order code"
// @Param driver_code query string false "Filter by driver code"
// @Param skip query int false "Number of items to skip" default(0)
// @Param limit query int false "Max items per page" default(20)
// @Success 200 {object} dto.PaginatedResponse
// @Failure 500 {object} map[string]interface{}
// @Router /tracking-logs [get]
func (h *Handler) List(c *gin.Context) {
	skip, _ := strconv.Atoi(c.DefaultQuery("skip", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	orderCode := c.Query("order_code")
	driverCode := c.Query("driver_code")
	if limit <= 0 {
		limit = 20
	}
	if skip < 0 {
		skip = 0
	}
	events, total, err := h.Service.List(c.Request.Context(), orderCode, driverCode, skip, limit)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.PaginatedResponse{
		Items: dto.ToResponseList(events),
		Total: total,
		Skip:  skip,
		Limit: limit,
	})
}

// @Summary List tracking events for an order
// @Tags Tracking
// @Produce json
// @Security BearerAuth
// @Param order_code path string true "Order code"
// @Success 200 {array} dto.TrackingEventResponse
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /tracking-logs/order/{order_code} [get]
func (h *Handler) GetByOrder(c *gin.Context) {
	orderCode := c.Param("order_code")
	events, err := h.Service.GetByOrder(c.Request.Context(), orderCode)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToResponseList(events))
}

// @Summary Get a tracking event by ID
// @Tags Tracking
// @Produce json
// @Security BearerAuth
// @Param id path int true "Tracking event ID"
// @Success 200 {object} dto.TrackingEventResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /tracking-logs/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(fmt.Errorf("%w: invalid id", apierrors.ErrBadRequest))
		return
	}
	event, err := h.Service.Get(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToResponse(event))
}

// @Summary Update a tracking event
// @Tags Tracking
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Tracking event ID"
// @Param event body dto.UpdateTrackingEventRequest true "Tracking event update"
// @Success 200 {object} dto.TrackingEventResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /tracking-logs/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(fmt.Errorf("%w: invalid id", apierrors.ErrBadRequest))
		return
	}
	var req dto.UpdateTrackingEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(fmt.Errorf("%w: %v", apierrors.ErrBadRequest, err))
		return
	}
	event, err := h.Service.Update(c.Request.Context(), id, &req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToResponse(event))
}

// @Summary Delete a tracking event
// @Tags Tracking
// @Security BearerAuth
// @Param id path int true "Tracking event ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /tracking-logs/{id} [delete]
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
