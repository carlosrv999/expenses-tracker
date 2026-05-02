package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"expenses/internal/model"
)

type ExpenseRepository struct {
	db *sql.DB
}

func NewExpenseRepository(db *sql.DB) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

type ExpenseFilter struct {
	CategoryID      *int64
	PaymentMethodID *int64
	StartDate       *time.Time
	EndDate         *time.Time
	TagIDs          []int64
	IncludeDeleted  bool
	Limit           int
	Offset          int

	// Relations allows the caller (usually the HTTP handler) to request optional
	// Many-to-One relations to be included in the response.
	// Supported values: "category", "payment_method"
	// Tags are *always* included (via JSON aggregation).
	// Example: ?relations=category,payment_method or ?relations=category&relations=payment_method
	// Backwards-compatible: if empty/nil, only expense + tags are returned (same as before).
	Relations []string
}

// contains is a small helper (add it near the ExpenseFilter type, e.g. after the struct).
// Case-insensitive to make query param handling more robust.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, item) {
			return true
		}
	}
	return false
}

func (r *ExpenseRepository) Create(ctx context.Context, e *model.Expense, tagIDs []int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	const insert = `
		INSERT INTO expense (category_id, payment_method_id, currency, amount, expense_date, merchant_name, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING expense_id, created_at, updated_at`
	if err := tx.QueryRowContext(ctx, insert,
		e.CategoryID, e.PaymentMethodID, e.Currency, e.Amount, e.ExpenseDate, e.MerchantName, e.Description,
	).Scan(&e.ExpenseID, &e.CreatedAt, &e.UpdatedAt); err != nil {
		return err
	}

	if err := setExpenseTagsTx(ctx, tx, e.ExpenseID, tagIDs); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *ExpenseRepository) GetByID(ctx context.Context, id int64) (*model.Expense, error) {
	const q = `
		SELECT expense_id, category_id, payment_method_id, currency, amount, expense_date,
		       merchant_name, description, created_at, updated_at, deleted_at
		FROM expense
		WHERE expense_id = $1 AND deleted_at IS NULL`
	var e model.Expense
	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&e.ExpenseID, &e.CategoryID, &e.PaymentMethodID, &e.Currency, &e.Amount, &e.ExpenseDate,
		&e.MerchantName, &e.Description, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &e, nil
}

type PaginatedExpenseList struct {
	Expenses   []model.Expense `json:"expenses"`
	TotalCount int64           `json:"total_count"`
	Limit      int             `json:"limit"`
	Offset     int             `json:"offset"`
}

// ListPaginated returns expenses with full pagination metadata (items + total count).
// It reuses the exact same filtering logic as List() for consistency.
// This is the recommended method when the frontend needs to show "Page X of Y", total records, etc.
func (r *ExpenseRepository) ListPaginated(ctx context.Context, f ExpenseFilter) (PaginatedExpenseList, error) {
	var (
		conds []string
		args  []any
	)

	// Count + data query now use the SAME alias 'e' everywhere → no more Postgres error
	if !f.IncludeDeleted {
		conds = append(conds, "e.deleted_at IS NULL")
	}
	if f.CategoryID != nil {
		args = append(args, *f.CategoryID)
		conds = append(conds, fmt.Sprintf("e.category_id = $%d", len(args)))
	}
	if f.PaymentMethodID != nil {
		args = append(args, *f.PaymentMethodID)
		conds = append(conds, fmt.Sprintf("e.payment_method_id = $%d", len(args)))
	}
	if f.StartDate != nil {
		args = append(args, *f.StartDate)
		conds = append(conds, fmt.Sprintf("e.expense_date >= $%d", len(args)))
	}
	if f.EndDate != nil {
		args = append(args, *f.EndDate)
		conds = append(conds, fmt.Sprintf("e.expense_date <= $%d", len(args)))
	}
	// Tag filter (ANY of the provided tags) – consistent alias 'e'
	if len(f.TagIDs) > 0 {
		tagPlaceholders := make([]string, len(f.TagIDs))
		for i, tagID := range f.TagIDs {
			args = append(args, tagID)
			tagPlaceholders[i] = fmt.Sprintf("$%d", len(args))
		}
		conds = append(conds, fmt.Sprintf(
			`EXISTS (SELECT 1 FROM expense_tag et WHERE et.expense_id = e.expense_id AND et.tag_id IN (%s))`,
			strings.Join(tagPlaceholders, ", "),
		))
	}

	where := ""
	if len(conds) > 0 {
		where = "WHERE " + strings.Join(conds, " AND ")
	}

	// 1. Get total count (now uses alias e)
	countQuery := fmt.Sprintf(`
        SELECT COUNT(*)
        FROM expense e
        %s`, where)

	var totalCount int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount); err != nil {
		return PaginatedExpenseList{}, err
	}

	// 2. Get the actual paginated rows — now with optional relations support
	limit := 100
	if f.Limit > 0 && f.Limit <= 500 {
		limit = f.Limit
	}

	// Copy filter args and append LIMIT/OFFSET
	dataArgs := make([]any, len(args))
	copy(dataArgs, args)
	dataArgs = append(dataArgs, limit)
	limitPlaceholder := fmt.Sprintf("$%d", len(dataArgs))
	dataArgs = append(dataArgs, f.Offset)
	offsetPlaceholder := fmt.Sprintf("$%d", len(dataArgs))

	// === SAME dynamic query building as in List() ===
	includeCategory := contains(f.Relations, "category")
	includePaymentMethod := contains(f.Relations, "payment_method")

	selectFields := []string{
		`e.expense_id, e.category_id, e.payment_method_id, e.currency, e.amount, 
		 e.expense_date, e.merchant_name, e.description, 
		 e.created_at, e.updated_at, e.deleted_at`,
	}
	joins := []string{
		`LEFT JOIN expense_tag et ON et.expense_id = e.expense_id`,
		`LEFT JOIN tag t ON t.tag_id = et.tag_id`,
	}
	groupByFields := []string{
		`e.expense_id, e.category_id, e.payment_method_id, e.currency, e.amount,
		 e.expense_date, e.merchant_name, e.description, 
		 e.created_at, e.updated_at, e.deleted_at`,
	}

	if includeCategory {
		selectFields = append(selectFields, `
			COALESCE(jsonb_build_object(
				'category_id',        c.category_id,
				'parent_category_id', c.parent_category_id,
				'category_name',      c.category_name,
				'icon',               c.icon,
				'color',              c.color
			), NULL::jsonb) AS category`)
		joins = append(joins, `LEFT JOIN category c ON c.category_id = e.category_id`)
		groupByFields = append(groupByFields, `c.category_id, c.parent_category_id, c.category_name, c.icon, c.color`)
	}

	if includePaymentMethod {
		selectFields = append(selectFields, `
			COALESCE(jsonb_build_object(
				'payment_method_id', pm.payment_method_id,
				'method_name',       pm.method_name,
				'method_type',       pm.method_type,
				'icon',              pm.icon
			), NULL::jsonb) AS payment_method`)
		joins = append(joins, `LEFT JOIN payment_method pm ON pm.payment_method_id = e.payment_method_id`)
		groupByFields = append(groupByFields, `pm.payment_method_id, pm.method_name, pm.method_type, pm.icon`)
	}

	// Tags always included
	selectFields = append(selectFields, `
		COALESCE(
			jsonb_agg(
				jsonb_build_object(
					'tag_id',     t.tag_id,
					'tag_name',   t.tag_name,
					'color',      t.color,
					'icon',       t.icon
				)
			) FILTER (WHERE t.tag_id IS NOT NULL),
			'[]'::jsonb
		) AS tags`)

	dataQuery := fmt.Sprintf(`
		SELECT 
			%s
		FROM expense e
		%s
		%s
		GROUP BY 
			%s
		ORDER BY e.expense_date DESC, e.expense_id DESC
		LIMIT %s OFFSET %s`,
		strings.Join(selectFields, ", "),
		strings.Join(joins, "\n\t\t"),
		where,
		strings.Join(groupByFields, ", "),
		limitPlaceholder,
		offsetPlaceholder)

	rows, err := r.db.QueryContext(ctx, dataQuery, dataArgs...)
	if err != nil {
		return PaginatedExpenseList{}, err
	}
	defer rows.Close()

	var expenses []model.Expense
	for rows.Next() {
		var e model.Expense
		var tagsJSON []byte
		var categoryJSON, paymentMethodJSON []byte

		// Dynamic scan (exact same logic as List)
		scanArgs := []interface{}{
			&e.ExpenseID, &e.CategoryID, &e.PaymentMethodID, &e.Currency, &e.Amount,
			&e.ExpenseDate, &e.MerchantName, &e.Description,
			&e.CreatedAt, &e.UpdatedAt, &e.DeletedAt,
		}
		if includeCategory {
			scanArgs = append(scanArgs, &categoryJSON)
		}
		if includePaymentMethod {
			scanArgs = append(scanArgs, &paymentMethodJSON)
		}
		scanArgs = append(scanArgs, &tagsJSON)

		if err := rows.Scan(scanArgs...); err != nil {
			return PaginatedExpenseList{}, err
		}

		// Unmarshal tags (always)
		if len(tagsJSON) > 0 {
			if err := json.Unmarshal(tagsJSON, &e.Tags); err != nil {
				return PaginatedExpenseList{}, fmt.Errorf("failed to unmarshal tags for expense %d: %w", e.ExpenseID, err)
			}
		}

		// Optional relations
		if includeCategory && len(categoryJSON) > 0 {
			var cat model.Category
			if err := json.Unmarshal(categoryJSON, &cat); err != nil {
				return PaginatedExpenseList{}, fmt.Errorf("failed to unmarshal category for expense %d: %w", e.ExpenseID, err)
			}
			e.Category = &cat
		}
		if includePaymentMethod && len(paymentMethodJSON) > 0 {
			var pm model.PaymentMethod
			if err := json.Unmarshal(paymentMethodJSON, &pm); err != nil {
				return PaginatedExpenseList{}, fmt.Errorf("failed to unmarshal payment method for expense %d: %w", e.ExpenseID, err)
			}
			e.PaymentMethod = &pm
		}

		expenses = append(expenses, e)
	}

	if err := rows.Err(); err != nil {
		return PaginatedExpenseList{}, err
	}

	return PaginatedExpenseList{
		Expenses:   expenses,
		TotalCount: totalCount,
		Limit:      limit,
		Offset:     f.Offset,
	}, nil
}

func (r *ExpenseRepository) Update(ctx context.Context, e *model.Expense, tagIDs *[]int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	const q = `
		UPDATE expense
		SET category_id       = $1,
		    payment_method_id = $2,
		    currency          = $3,
		    amount            = $4,
		    expense_date      = $5,
		    merchant_name     = $6,
		    description       = $7,
		    updated_at        = NOW()
		WHERE expense_id = $8 AND deleted_at IS NULL
		RETURNING updated_at`
	err = tx.QueryRowContext(ctx, q,
		e.CategoryID, e.PaymentMethodID, e.Currency, e.Amount, e.ExpenseDate,
		e.MerchantName, e.Description, e.ExpenseID,
	).Scan(&e.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}

	if tagIDs != nil {
		if err := setExpenseTagsTx(ctx, tx, e.ExpenseID, *tagIDs); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *ExpenseRepository) SoftDelete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE expense SET deleted_at = NOW(), updated_at = NOW() WHERE expense_id = $1 AND deleted_at IS NULL`, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func setExpenseTagsTx(ctx context.Context, tx *sql.Tx, expenseID int64, tagIDs []int64) error {
	if _, err := tx.ExecContext(ctx, `DELETE FROM expense_tag WHERE expense_id = $1`, expenseID); err != nil {
		return err
	}
	if len(tagIDs) == 0 {
		return nil
	}

	values := make([]string, 0, len(tagIDs))
	args := make([]any, 0, len(tagIDs)*2)
	args = append(args, expenseID)
	for i, tagID := range tagIDs {
		values = append(values, fmt.Sprintf("($1, $%d)", i+2))
		args = append(args, tagID)
	}
	q := fmt.Sprintf(`INSERT INTO expense_tag (expense_id, tag_id) VALUES %s ON CONFLICT DO NOTHING`,
		strings.Join(values, ", "))
	_, err := tx.ExecContext(ctx, q, args...)
	return err
}

func (r *ExpenseRepository) BulkCreate(ctx context.Context, expenses []*model.Expense, tagIDsList [][]int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	const insert = `
		INSERT INTO expense (category_id, payment_method_id, currency, amount, expense_date, merchant_name, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING expense_id, created_at, updated_at`

	for i, e := range expenses {
		if err := tx.QueryRowContext(ctx, insert,
			e.CategoryID, e.PaymentMethodID, e.Currency, e.Amount, e.ExpenseDate, e.MerchantName, e.Description,
		).Scan(&e.ExpenseID, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return fmt.Errorf("failed to insert expense %d: %w", i, err)
		}

		if err := setExpenseTagsTx(ctx, tx, e.ExpenseID, tagIDsList[i]); err != nil {
			return fmt.Errorf("failed to set tags for expense %d: %w", i, err)
		}
	}

	return tx.Commit()
}
