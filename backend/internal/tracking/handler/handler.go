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

func (h *Handler) GetByOrder(c *gin.Context) {
	orderCode := c.Param("order_code")
	events, err := h.Service.GetByOrder(c.Request.Context(), orderCode)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToResponseList(events))
}

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
