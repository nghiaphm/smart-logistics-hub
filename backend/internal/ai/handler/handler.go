package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"my-web-app.com/smart-logistic-hub/internal/ai/dto"
	"my-web-app.com/smart-logistic-hub/internal/ai/service"
	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
)

type Handler struct {
	Service *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{Service: svc}
}

// @Summary Create an AI event
// @Tags AI Events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param event body dto.CreateAIEventRequest true "AI event payload"
// @Success 201 {object} dto.AIEventResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /ai-events [post]
func (h *Handler) Create(c *gin.Context) {
	var req dto.CreateAIEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(fmt.Errorf("%w: %v", apierrors.ErrBadRequest, err))
		return
	}

	e, warning, err := h.Service.Create(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}
	resp := dto.ToResponse(e, warning != "")
	c.JSON(http.StatusCreated, resp)
}

// @Summary List AI events
// @Tags AI Events
// @Produce json
// @Security BearerAuth
// @Param license_plate query string false "Filter by license plate"
// @Param gate_id query string false "Filter by gate ID"
// @Param event_type query string false "Filter by event type"
// @Param skip query int false "Number of items to skip" default(0)
// @Param limit query int false "Max items per page" default(20)
// @Success 200 {object} dto.PaginatedResponse
// @Failure 500 {object} map[string]interface{}
// @Router /ai-events [get]
func (h *Handler) List(c *gin.Context) {
	skip, _ := strconv.Atoi(c.DefaultQuery("skip", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit <= 0 {
		limit = 20
	}
	if skip < 0 {
		skip = 0
	}

	items, total, err := h.Service.List(c.Request.Context(),
		c.Query("license_plate"), c.Query("gate_id"), c.Query("event_type"), skip, limit)
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

// @Summary Get an AI event by ID
// @Tags AI Events
// @Produce json
// @Security BearerAuth
// @Param id path int true "AI event ID"
// @Success 200 {object} dto.AIEventResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /ai-events/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(fmt.Errorf("%w: invalid id", apierrors.ErrBadRequest))
		return
	}

	e, err := h.Service.Get(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToResponse(e, e.ConfidenceScore < 0.7))
}

// @Summary Update an AI event
// @Tags AI Events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "AI event ID"
// @Param event body dto.UpdateAIEventRequest true "AI event update"
// @Success 200 {object} dto.AIEventResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /ai-events/{id} [patch]
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(fmt.Errorf("%w: invalid id", apierrors.ErrBadRequest))
		return
	}

	var req dto.UpdateAIEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(fmt.Errorf("%w: %v", apierrors.ErrBadRequest, err))
		return
	}

	e, err := h.Service.Update(c.Request.Context(), id, &req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToResponse(e, e.ConfidenceScore < 0.7))
}

// @Summary Delete an AI event
// @Tags AI Events
// @Security BearerAuth
// @Param id path int true "AI event ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /ai-events/{id} [delete]
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
