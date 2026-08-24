package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	"my-web-app.com/smart-logistic-hub/internal/driver/dto"
	"my-web-app.com/smart-logistic-hub/internal/driver/service"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

// @Summary Create a driver
// @Tags Drivers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param driver body dto.CreateDriverRequest true "Driver payload"
// @Success 201 {object} dto.DriverResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /drivers [post]
func (h *Handler) Create(c *gin.Context) {
	var req dto.CreateDriverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(fmt.Errorf("%w: %v", apierrors.ErrBadRequest, err))
		return
	}
	d, err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, dto.ToResponse(d))
}

// @Summary Get a driver by ID
// @Tags Drivers
// @Produce json
// @Security BearerAuth
// @Param id path int true "Driver ID"
// @Success 200 {object} dto.DriverResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /drivers/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(fmt.Errorf("%w: invalid driver ID", apierrors.ErrBadRequest))
		return
	}
	d, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToResponse(d))
}

// @Summary List drivers
// @Tags Drivers
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filter by driver status"
// @Param skip query int false "Number of items to skip" default(0)
// @Param limit query int false "Max items per page" default(20)
// @Success 200 {object} dto.PaginatedResponse
// @Failure 500 {object} map[string]interface{}
// @Router /drivers [get]
func (h *Handler) List(c *gin.Context) {
	status := c.Query("status")
	offset, _ := strconv.Atoi(c.DefaultQuery("skip", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit < 1 || limit > 200 {
		limit = 20
	}

	drivers, total, err := h.svc.List(c.Request.Context(), status, offset, limit)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.PaginatedResponse{
		Items: dto.ToResponseList(drivers),
		Total: total,
		Skip:  offset,
		Limit: limit,
	})
}

// @Summary Update a driver
// @Tags Drivers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Driver ID"
// @Param driver body dto.UpdateDriverRequest true "Driver update"
// @Success 200 {object} dto.DriverResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /drivers/{id} [patch]
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(fmt.Errorf("%w: invalid driver ID", apierrors.ErrBadRequest))
		return
	}
	var req dto.UpdateDriverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(fmt.Errorf("%w: %v", apierrors.ErrBadRequest, err))
		return
	}
	d, err := h.svc.Update(c.Request.Context(), id, &req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToResponse(d))
}

// @Summary Delete a driver
// @Tags Drivers
// @Security BearerAuth
// @Param id path int true "Driver ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /drivers/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(fmt.Errorf("%w: invalid driver ID", apierrors.ErrBadRequest))
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}
