package category

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"

	"expenses/internal/handler"
	"expenses/internal/model"
	"expenses/internal/service"
)

type Handler struct {
	svc *service.CategoryService
}

func NewHandler(svc *service.CategoryService) *Handler {
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
	if strings.TrimSpace(req.CategoryName) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "category_name is required")
	}

	cat := &model.Category{
		ParentCategoryID: req.ParentCategoryID,
		CategoryName:     req.CategoryName,
		Icon:             req.Icon,
		Color:            req.Color,
	}
	if err := h.svc.Create(c.Request().Context(), cat); err != nil {
		return handler.MapError(err)
	}
	return c.JSON(http.StatusCreated, cat)
}

func (h *Handler) get(c *echo.Context) error {
	id, err := handler.ParseIDParam(c, "id")
	if err != nil {
		return err
	}
	cat, err := h.svc.Get(c.Request().Context(), id)
	if err != nil {
		return handler.MapError(err)
	}
	return c.JSON(http.StatusOK, cat)
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
	if strings.TrimSpace(req.CategoryName) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "category_name is required")
	}
	cat := &model.Category{
		CategoryID:       id,
		ParentCategoryID: req.ParentCategoryID,
		CategoryName:     req.CategoryName,
		Icon:             req.Icon,
		Color:            req.Color,
	}
	if err := h.svc.Update(c.Request().Context(), cat); err != nil {
		return handler.MapError(err)
	}
	return c.JSON(http.StatusOK, cat)
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
