package paymentmethod

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"

	"expenses/internal/handler"
	"expenses/internal/model"
	"expenses/internal/service"
)

type Handler struct {
	svc *service.PaymentMethodService
}

func NewHandler(svc *service.PaymentMethodService) *Handler {
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
	if strings.TrimSpace(req.MethodName) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "method_name is required")
	}

	pm := &model.PaymentMethod{
		MethodName: req.MethodName,
		MethodType: model.PaymentMethodType(req.MethodType),
		Icon:       req.Icon,
	}
	if err := h.svc.Create(c.Request().Context(), pm); err != nil {
		return mapServiceError(err)
	}
	return c.JSON(http.StatusCreated, pm)
}

func (h *Handler) get(c *echo.Context) error {
	id, err := handler.ParseIDParam(c, "id")
	if err != nil {
		return err
	}
	pm, err := h.svc.Get(c.Request().Context(), id)
	if err != nil {
		return handler.MapError(err)
	}
	return c.JSON(http.StatusOK, pm)
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
	if strings.TrimSpace(req.MethodName) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "method_name is required")
	}
	pm := &model.PaymentMethod{
		PaymentMethodID: id,
		MethodName:      req.MethodName,
		MethodType:      model.PaymentMethodType(req.MethodType),
		Icon:            req.Icon,
	}
	if err := h.svc.Update(c.Request().Context(), pm); err != nil {
		return mapServiceError(err)
	}
	return c.JSON(http.StatusOK, pm)
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

func mapServiceError(err error) error {
	if errors.Is(err, service.ErrInvalidPaymentMethodType) {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return handler.MapError(err)
}
