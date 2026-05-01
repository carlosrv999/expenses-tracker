package expense

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

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
	g.GET("", h.list) // ← now returns paginated result
	g.POST("", h.create)
	g.GET("/:id", h.get)
	g.PUT("/:id", h.update)
	g.DELETE("/:id", h.delete)

	// ✅ CSV bulk upload
	g.POST("/upload", h.uploadCSV)
}

// Updated: now uses ListPaginated (full pagination metadata + tags attached)
func (h *Handler) list(c *echo.Context) error {
	f, err := parseFilter(c)
	if err != nil {
		return err
	}

	result, err := h.svc.ListPaginated(c.Request().Context(), f)
	if err != nil {
		return handler.MapError(err)
	}

	return c.JSON(http.StatusOK, result)
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

	// === Date range filters ===
	if v := c.QueryParam("start_date"); v != "" {
		t, err := parseDate(v)
		if err != nil {
			return f, echo.NewHTTPError(http.StatusBadRequest, "invalid start_date. Use format YYYY-MM-DD (or RFC3339)")
		}
		f.StartDate = &t
	}
	if v := c.QueryParam("end_date"); v != "" {
		t, err := parseDate(v)
		if err != nil {
			return f, echo.NewHTTPError(http.StatusBadRequest, "invalid end_date. Use format YYYY-MM-DD (or RFC3339)")
		}
		// Make end_date inclusive of the entire day
		year, month, day := t.Date()
		endOfDay := time.Date(year, month, day, 23, 59, 59, 999999999, t.Location())
		f.EndDate = &endOfDay
	}

	// === NEW: Tags filter (supports ?tags=[1,3], ?tags=1,3 or ?tags=1&tags=3) ===
	if tagValues := c.QueryParams()["tags"]; len(tagValues) > 0 {
		tagStr := strings.Join(tagValues, ",")
		tagIDs, err := parseTagIDs(tagStr)
		if err != nil {
			return f, echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid tags parameter: %v", err))
		}
		f.TagIDs = tagIDs
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

// Upload CSV file
func (h *Handler) uploadCSV(c *echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "file is required (field name: 'file')")
	}

	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "cannot open file")
	}
	defer src.Close()

	expenses, tagIDsList, err := parseExpenseCSV(src)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("CSV parsing error: %v", err))
	}

	if len(expenses) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "CSV file is empty")
	}

	if err := h.svc.BulkCreate(c.Request().Context(), expenses, tagIDsList); err != nil {
		return mapServiceError(err)
	}

	return c.JSON(http.StatusCreated, map[string]any{
		"message":  "expenses uploaded successfully",
		"uploaded": len(expenses),
	})
}

// ... (parseExpenseCSV, parseDate, parseTagIDs remain unchanged - I kept them exactly as you had)
func parseExpenseCSV(r io.Reader) ([]*model.Expense, [][]int64, error) {
	// [your existing parseExpenseCSV function - no changes needed]
	reader := csv.NewReader(r)
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, fmt.Errorf("invalid CSV: %w", err)
	}
	if len(records) < 2 {
		return nil, nil, errors.New("CSV must have header + at least one data row")
	}

	var expenses []*model.Expense
	var tagIDsList [][]int64

	for i := 1; i < len(records); i++ {
		row := records[i]
		if len(row) < 4 {
			return nil, nil, fmt.Errorf("row %d: not enough columns", i)
		}

		expenseDate, err := parseDate(row[0])
		if err != nil {
			return nil, nil, fmt.Errorf("row %d: %w", i, err)
		}

		amount, err := strconv.ParseInt(strings.TrimSpace(row[1]), 10, 64)
		if err != nil || amount <= 0 {
			return nil, nil, fmt.Errorf("row %d: invalid amount", i)
		}

		categoryID, _ := strconv.ParseInt(strings.TrimSpace(row[2]), 10, 64)
		paymentMethodID, _ := strconv.ParseInt(strings.TrimSpace(row[3]), 10, 64)

		currency := "PEN"
		if len(row) > 4 && strings.TrimSpace(row[4]) != "" {
			currency = strings.TrimSpace(row[4])
		}

		var merchantName *string
		if len(row) > 5 && strings.TrimSpace(row[5]) != "" {
			s := strings.TrimSpace(row[5])
			merchantName = &s
		}

		var description *string
		if len(row) > 6 && strings.TrimSpace(row[6]) != "" {
			s := strings.TrimSpace(row[6])
			description = &s
		}

		var tagIDs []int64
		if len(row) > 7 {
			tagIDs, err = parseTagIDs(row[7])
			if err != nil {
				return nil, nil, fmt.Errorf("row %d: %w", i, err)
			}
		}

		e := &model.Expense{
			CategoryID:      categoryID,
			PaymentMethodID: paymentMethodID,
			Currency:        currency,
			Amount:          amount,
			ExpenseDate:     expenseDate,
			MerchantName:    merchantName,
			Description:     description,
		}

		expenses = append(expenses, e)
		tagIDsList = append(tagIDsList, tagIDs)
	}

	return expenses, tagIDsList, nil
}

func parseDate(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, errors.New("expense_date is required")
	}
	layouts := []string{time.DateOnly, time.RFC3339, "2006-01-02 15:04:05"}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid date format: %s (use YYYY-MM-DD)", s)
}

func parseTagIDs(s string) ([]int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}

	// Support both "1,2,3" and "[1,2,3]" formats
	if len(s) > 1 && s[0] == '[' && s[len(s)-1] == ']' {
		s = s[1 : len(s)-1]
	}

	parts := strings.Split(s, ",")
	ids := make([]int64, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		id, err := strconv.ParseInt(p, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid tag_id '%s'", p)
		}
		ids = append(ids, id)
	}
	return ids, nil
}
