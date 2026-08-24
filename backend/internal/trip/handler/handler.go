package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	"my-web-app.com/smart-logistic-hub/internal/trip/dto"
	"my-web-app.com/smart-logistic-hub/internal/trip/service"
)

type Handler struct {
	Service *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{Service: svc}
}

// @Summary Create a trip
// @Tags Trips
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param trip body dto.CreateTripRequest true "Trip payload"
// @Success 201 {object} dto.TripResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /trips [post]
func (h *Handler) Create(c *gin.Context) {
	var req dto.CreateTripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(fmt.Errorf("%w: %v", apierrors.ErrBadRequest, err))
		return
	}

	t, err := h.Service.Create(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, dto.ToResponse(t, nil))
}

// @Summary List trips
// @Tags Trips
// @Produce json
// @Security BearerAuth
// @Param skip query int false "Number of items to skip" default(0)
// @Param limit query int false "Max items per page" default(20)
// @Success 200 {object} dto.PaginatedResponse
// @Failure 500 {object} map[string]interface{}
// @Router /trips [get]
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

// @Summary Get a trip by ID
// @Tags Trips
// @Produce json
// @Security BearerAuth
// @Param id path int true "Trip ID"
// @Success 200 {object} dto.TripResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /trips/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(fmt.Errorf("%w: invalid id", apierrors.ErrBadRequest))
		return
	}

	t, stops, err := h.Service.Get(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToResponse(t, stops))
}

// @Summary Update a trip
// @Tags Trips
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Trip ID"
// @Param trip body dto.UpdateTripRequest true "Trip update"
// @Success 200 {object} dto.TripResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /trips/{id} [patch]
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(fmt.Errorf("%w: invalid id", apierrors.ErrBadRequest))
		return
	}

	var req dto.UpdateTripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(fmt.Errorf("%w: %v", apierrors.ErrBadRequest, err))
		return
	}

	t, err := h.Service.Update(c.Request.Context(), id, &req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToResponse(t, nil))
}

// @Summary Assign a driver to a trip
// @Tags Trips
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Trip ID"
// @Param request body dto.AssignDriverRequest true "Driver code"
// @Success 200 {object} dto.TripResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /trips/{id}/assign-driver [post]
func (h *Handler) AssignDriver(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(fmt.Errorf("%w: invalid id", apierrors.ErrBadRequest))
		return
	}

	var req dto.AssignDriverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(fmt.Errorf("%w: %v", apierrors.ErrBadRequest, err))
		return
	}

	t, err := h.Service.AssignDriver(c.Request.Context(), id, req.DriverCode)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToResponse(t, nil))
}

// @Summary Delete a trip
// @Tags Trips
// @Security BearerAuth
// @Param id path int true "Trip ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /trips/{id} [delete]
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
