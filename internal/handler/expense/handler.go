package expense

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"

	"expenses/internal/handler"
	"expenses/internal/model"
	"expenses/internal/repository"
	"expenses/internal/service"
)

type Handler struct {
	svc *service.ExpenseService
}

func NewHandler(svc *service.ExpenseService) *Handler {
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
	f, err := parseFilter(c)
	if err != nil {
		return err
	}
	items, err := h.svc.List(c.Request().Context(), f)
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
	currency := req.Currency
	if currency == "" {
		currency = "PEN"
	}
	e := &model.Expense{
		CategoryID:      req.CategoryID,
		PaymentMethodID: req.PaymentMethodID,
		Currency:        currency,
		Amount:          req.Amount,
		ExpenseDate:     req.ExpenseDate,
		MerchantName:    req.MerchantName,
		Description:     req.Description,
	}
	if err := h.svc.Create(c.Request().Context(), e, req.TagIDs); err != nil {
		return mapServiceError(err)
	}
	return c.JSON(http.StatusCreated, e)
}

func (h *Handler) get(c *echo.Context) error {
	id, err := handler.ParseIDParam(c, "id")
	if err != nil {
		return err
	}
	e, err := h.svc.Get(c.Request().Context(), id)
	if err != nil {
		return handler.MapError(err)
	}
	return c.JSON(http.StatusOK, e)
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
	currency := req.Currency
	if currency == "" {
		currency = "PEN"
	}
	e := &model.Expense{
		ExpenseID:       id,
		CategoryID:      req.CategoryID,
		PaymentMethodID: req.PaymentMethodID,
		Currency:        currency,
		Amount:          req.Amount,
		ExpenseDate:     req.ExpenseDate,
		MerchantName:    req.MerchantName,
		Description:     req.Description,
	}
	if err := h.svc.Update(c.Request().Context(), e, req.TagIDs); err != nil {
		return mapServiceError(err)
	}
	return c.JSON(http.StatusOK, e)
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

func parseFilter(c *echo.Context) (repository.ExpenseFilter, error) {
	var f repository.ExpenseFilter

	if v := c.QueryParam("category_id"); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return f, echo.NewHTTPError(http.StatusBadRequest, "invalid category_id")
		}
		f.CategoryID = &id
	}
	if v := c.QueryParam("payment_method_id"); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return f, echo.NewHTTPError(http.StatusBadRequest, "invalid payment_method_id")
		}
		f.PaymentMethodID = &id
	}
	if v := c.QueryParam("limit"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 0 {
			return f, echo.NewHTTPError(http.StatusBadRequest, "invalid limit")
		}
		f.Limit = n
	}
	if v := c.QueryParam("offset"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 0 {
			return f, echo.NewHTTPError(http.StatusBadRequest, "invalid offset")
		}
		f.Offset = n
	}
	return f, nil
}

func mapServiceError(err error) error {
	if errors.Is(err, service.ErrInvalidAmount) || errors.Is(err, service.ErrInvalidCurrency) {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return handler.MapError(err)
}
