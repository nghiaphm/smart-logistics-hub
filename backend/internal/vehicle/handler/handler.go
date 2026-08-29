package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	"my-web-app.com/smart-logistic-hub/internal/vehicle/dto"
	"my-web-app.com/smart-logistic-hub/internal/vehicle/service"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

// @Summary Create a vehicle
// @Tags Vehicles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param vehicle body dto.CreateVehicleRequest true "Vehicle payload"
// @Success 201 {object} dto.VehicleResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /vehicles [post]
func (h *Handler) Create(c *gin.Context) {
	var req dto.CreateVehicleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(fmt.Errorf("%w: %v", apierrors.ErrBadRequest, err))
		return
	}
	v, err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, dto.ToResponse(v))
}

// @Summary Get a vehicle by ID
// @Tags Vehicles
// @Produce json
// @Security BearerAuth
// @Param id path int true "Vehicle ID"
// @Success 200 {object} dto.VehicleResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /vehicles/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(fmt.Errorf("%w: invalid vehicle ID", apierrors.ErrBadRequest))
		return
	}
	v, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToResponse(v))
}

// @Summary List vehicles
// @Tags Vehicles
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filter by vehicle status"
// @Param skip query int false "Number of items to skip" default(0)
// @Param limit query int false "Max items per page" default(20)
// @Success 200 {object} dto.PaginatedResponse
// @Failure 500 {object} map[string]interface{}
// @Router /vehicles [get]
func (h *Handler) List(c *gin.Context) {
	status := c.Query("status")
	offset, _ := strconv.Atoi(c.DefaultQuery("skip", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit < 1 || limit > 200 {
		limit = 20
	}

	vehicles, total, err := h.svc.List(c.Request.Context(), status, offset, limit)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.PaginatedResponse{
		Items: dto.ToResponseList(vehicles),
		Total: total,
		Skip:  offset,
		Limit: limit,
	})
}

// @Summary Update a vehicle
// @Tags Vehicles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Vehicle ID"
// @Param vehicle body dto.UpdateVehicleRequest true "Vehicle update"
// @Success 200 {object} dto.VehicleResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /vehicles/{id} [patch]
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(fmt.Errorf("%w: invalid vehicle ID", apierrors.ErrBadRequest))
		return
	}
	var req dto.UpdateVehicleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(fmt.Errorf("%w: %v", apierrors.ErrBadRequest, err))
		return
	}
	v, err := h.svc.Update(c.Request.Context(), id, &req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToResponse(v))
}

// @Summary Delete a vehicle
// @Tags Vehicles
// @Security BearerAuth
// @Param id path int true "Vehicle ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /vehicles/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(fmt.Errorf("%w: invalid vehicle ID", apierrors.ErrBadRequest))
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}
