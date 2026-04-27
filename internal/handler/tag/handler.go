package tag

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"

	"expenses/internal/handler"
	"expenses/internal/model"
	"expenses/internal/service"
)

type Handler struct {
	svc *service.TagService
}

func NewHandler(svc *service.TagService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(g *echo.Group) {
	g.GET("", h.list)
	g.POST("", h.create)
	g.GET("/:id", h.get)
	g.PUT("/:id", h.update)
	g.DELETE("/:id", h.delete)
}

func (h *Handler) list(c *echo.Context) error {
	items, err := h.svc.List(c.Request().Context())
	if err != nil {
		return handler.MapError(err)
	}
	return c.JSON(http.StatusOK, items)
}

func (h *Handler) create(c *echo.Context) error {
	var req createRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}
	if strings.TrimSpace(req.TagName) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "tag_name is required")
	}
	t := &model.Tag{
		TagName: req.TagName,
		Color:   req.Color,
		Icon:    req.Icon,
	}
	if err := h.svc.Create(c.Request().Context(), t); err != nil {
		return handler.MapError(err)
	}
	return c.JSON(http.StatusCreated, t)
}

func (h *Handler) get(c *echo.Context) error {
	id, err := handler.ParseIDParam(c, "id")
	if err != nil {
		return err
	}
	t, err := h.svc.Get(c.Request().Context(), id)
	if err != nil {
		return handler.MapError(err)
	}
	return c.JSON(http.StatusOK, t)
}

func (h *Handler) update(c *echo.Context) error {
	id, err := handler.ParseIDParam(c, "id")
	if err != nil {
		return err
	}
	var req updateRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}
	if strings.TrimSpace(req.TagName) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "tag_name is required")
	}
	t := &model.Tag{
		TagID:   id,
		TagName: req.TagName,
		Color:   req.Color,
		Icon:    req.Icon,
	}
	if err := h.svc.Update(c.Request().Context(), t); err != nil {
		return handler.MapError(err)
	}
	return c.JSON(http.StatusOK, t)
}

func (h *Handler) delete(c *echo.Context) error {
	id, err := handler.ParseIDParam(c, "id")
	if err != nil {
		return err
	}
	if err := h.svc.Delete(c.Request().Context(), id); err != nil {
		return handler.MapError(err)
	}
	return c.NoContent(http.StatusNoContent)
}
