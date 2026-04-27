package model

import "time"

type Category struct {
	CategoryID       int64     `json:"category_id"`
	ParentCategoryID *int64    `json:"parent_category_id,omitempty"`
	CategoryName     string    `json:"category_name"`
	Icon             *string   `json:"icon,omitempty"`
	Color            *string   `json:"color,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
