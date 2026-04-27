package repository

import (
	"context"
	"database/sql"
	"errors"

	"expenses/internal/model"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(ctx context.Context, c *model.Category) error {
	const q = `
		INSERT INTO category (parent_category_id, category_name, icon, color)
		VALUES ($1, $2, $3, $4)
		RETURNING category_id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, q, c.ParentCategoryID, c.CategoryName, c.Icon, c.Color).
		Scan(&c.CategoryID, &c.CreatedAt, &c.UpdatedAt)
}

func (r *CategoryRepository) GetByID(ctx context.Context, id int64) (*model.Category, error) {
	const q = `
		SELECT category_id, parent_category_id, category_name, icon, color, created_at, updated_at
		FROM category
		WHERE category_id = $1`
	var c model.Category
	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&c.CategoryID, &c.ParentCategoryID, &c.CategoryName, &c.Icon, &c.Color, &c.CreatedAt, &c.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepository) List(ctx context.Context) ([]model.Category, error) {
	const q = `
		SELECT category_id, parent_category_id, category_name, icon, color, created_at, updated_at
		FROM category
		ORDER BY category_id`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.Category
	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.CategoryID, &c.ParentCategoryID, &c.CategoryName, &c.Icon, &c.Color, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (r *CategoryRepository) Update(ctx context.Context, c *model.Category) error {
	const q = `
		UPDATE category
		SET parent_category_id = $1,
		    category_name      = $2,
		    icon               = $3,
		    color              = $4,
		    updated_at         = NOW()
		WHERE category_id = $5
		RETURNING updated_at`
	err := r.db.QueryRowContext(ctx, q, c.ParentCategoryID, c.CategoryName, c.Icon, c.Color, c.CategoryID).
		Scan(&c.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	return err
}

func (r *CategoryRepository) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM category WHERE category_id = $1`, id)
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
