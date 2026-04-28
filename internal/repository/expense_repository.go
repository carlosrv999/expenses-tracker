package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

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
	IncludeDeleted  bool
	Limit           int
	Offset          int
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

func (r *ExpenseRepository) List(ctx context.Context, f ExpenseFilter) ([]model.Expense, error) {
	var (
		conds []string
		args  []any
	)

	if !f.IncludeDeleted {
		conds = append(conds, "deleted_at IS NULL")
	}
	if f.CategoryID != nil {
		args = append(args, *f.CategoryID)
		conds = append(conds, fmt.Sprintf("category_id = $%d", len(args)))
	}
	if f.PaymentMethodID != nil {
		args = append(args, *f.PaymentMethodID)
		conds = append(conds, fmt.Sprintf("payment_method_id = $%d", len(args)))
	}

	where := ""
	if len(conds) > 0 {
		where = "WHERE " + strings.Join(conds, " AND ")
	}

	limit := 100
	if f.Limit > 0 && f.Limit <= 500 {
		limit = f.Limit
	}
	args = append(args, limit)
	limitPlaceholder := fmt.Sprintf("$%d", len(args))

	args = append(args, f.Offset)
	offsetPlaceholder := fmt.Sprintf("$%d", len(args))

	q := fmt.Sprintf(`
		SELECT expense_id, category_id, payment_method_id, currency, amount, expense_date,
		       merchant_name, description, created_at, updated_at, deleted_at
		FROM expense
		%s
		ORDER BY expense_date DESC, expense_id DESC
		LIMIT %s OFFSET %s`, where, limitPlaceholder, offsetPlaceholder)

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.Expense
	for rows.Next() {
		var e model.Expense
		if err := rows.Scan(
			&e.ExpenseID, &e.CategoryID, &e.PaymentMethodID, &e.Currency, &e.Amount, &e.ExpenseDate,
			&e.MerchantName, &e.Description, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
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
