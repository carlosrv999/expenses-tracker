package repository

import (
	"context"
	"database/sql"
	"errors"

	"expenses/internal/model"
)

type PaymentMethodRepository struct {
	db *sql.DB
}

func NewPaymentMethodRepository(db *sql.DB) *PaymentMethodRepository {
	return &PaymentMethodRepository{db: db}
}

func (r *PaymentMethodRepository) Create(ctx context.Context, p *model.PaymentMethod) error {
	const q = `
		INSERT INTO payment_method (method_name, method_type, icon)
		VALUES ($1, $2, $3)
		RETURNING payment_method_id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, q, p.MethodName, p.MethodType, p.Icon).
		Scan(&p.PaymentMethodID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *PaymentMethodRepository) GetByID(ctx context.Context, id int64) (*model.PaymentMethod, error) {
	const q = `
		SELECT payment_method_id, method_name, method_type, icon, created_at, updated_at
		FROM payment_method
		WHERE payment_method_id = $1`
	var p model.PaymentMethod
	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&p.PaymentMethodID, &p.MethodName, &p.MethodType, &p.Icon, &p.CreatedAt, &p.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PaymentMethodRepository) List(ctx context.Context) ([]model.PaymentMethod, error) {
	const q = `
		SELECT payment_method_id, method_name, method_type, icon, created_at, updated_at
		FROM payment_method
		ORDER BY payment_method_id`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.PaymentMethod
	for rows.Next() {
		var p model.PaymentMethod
		if err := rows.Scan(&p.PaymentMethodID, &p.MethodName, &p.MethodType, &p.Icon, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (r *PaymentMethodRepository) Update(ctx context.Context, p *model.PaymentMethod) error {
	const q = `
		UPDATE payment_method
		SET method_name = $1,
		    method_type = $2,
		    icon        = $3,
		    updated_at  = NOW()
		WHERE payment_method_id = $4
		RETURNING updated_at`
	err := r.db.QueryRowContext(ctx, q, p.MethodName, p.MethodType, p.Icon, p.PaymentMethodID).
		Scan(&p.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	return err
}

func (r *PaymentMethodRepository) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM payment_method WHERE payment_method_id = $1`, id)
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
