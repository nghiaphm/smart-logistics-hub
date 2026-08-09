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

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup, authMw gin.HandlerFunc) {
	g := rg.Group("/drivers")
	g.Use(authMw)

	g.POST("", h.create)
	g.GET("", h.list)
	g.GET("/:id", h.get)
	g.PUT("/:id", h.update)
	g.DELETE("/:id", h.delete)
}

func (h *Handler) create(c *gin.Context) {
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
	c.JSON(http.StatusCreated, d)
}

func (h *Handler) get(c *gin.Context) {
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
	c.JSON(http.StatusOK, d)
}

func (h *Handler) list(c *gin.Context) {
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
	c.JSON(http.StatusOK, gin.H{
		"items": drivers,
		"total": total,
		"skip":  offset,
		"limit": limit,
	})
}

func (h *Handler) update(c *gin.Context) {
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
	c.JSON(http.StatusOK, d)
}

func (h *Handler) delete(c *gin.Context) {
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
