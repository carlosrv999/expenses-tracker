package repository

import (
	"context"
	"database/sql"
	"errors"

	"expenses/internal/model"
)

type TagRepository struct {
	db *sql.DB
}

func NewTagRepository(db *sql.DB) *TagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) Create(ctx context.Context, t *model.Tag) error {
	const q = `
		INSERT INTO tag (tag_name, color, icon)
		VALUES ($1, $2, $3)
		RETURNING tag_id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, q, t.TagName, t.Color, t.Icon).
		Scan(&t.TagID, &t.CreatedAt, &t.UpdatedAt)
}

func (r *TagRepository) GetByID(ctx context.Context, id int64) (*model.Tag, error) {
	const q = `
		SELECT tag_id, tag_name, color, icon, created_at, updated_at
		FROM tag
		WHERE tag_id = $1`
	var t model.Tag
	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&t.TagID, &t.TagName, &t.Color, &t.Icon, &t.CreatedAt, &t.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TagRepository) List(ctx context.Context) ([]model.Tag, error) {
	const q = `
		SELECT tag_id, tag_name, color, icon, created_at, updated_at
		FROM tag
		ORDER BY tag_id`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.Tag
	for rows.Next() {
		var t model.Tag
		if err := rows.Scan(&t.TagID, &t.TagName, &t.Color, &t.Icon, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (r *TagRepository) Update(ctx context.Context, t *model.Tag) error {
	const q = `
		UPDATE tag
		SET tag_name   = $1,
		    color      = $2,
		    icon       = $3,
		    updated_at = NOW()
		WHERE tag_id = $4
		RETURNING updated_at`
	err := r.db.QueryRowContext(ctx, q, t.TagName, t.Color, t.Icon, t.TagID).
		Scan(&t.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	return err
}

func (r *TagRepository) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM tag WHERE tag_id = $1`, id)
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

func (r *TagRepository) ListByExpenseID(ctx context.Context, expenseID int64) ([]model.Tag, error) {
	const q = `
		SELECT t.tag_id, t.tag_name, t.color, t.icon, t.created_at, t.updated_at
		FROM tag t
		JOIN expense_tag et ON et.tag_id = t.tag_id
		WHERE et.expense_id = $1
		ORDER BY t.tag_id`
	rows, err := r.db.QueryContext(ctx, q, expenseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.Tag
	for rows.Next() {
		var t model.Tag
		if err := rows.Scan(&t.TagID, &t.TagName, &t.Color, &t.Icon, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}
