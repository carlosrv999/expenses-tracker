package model

import "time"

type Tag struct {
	TagID     int64     `json:"tag_id"`
	TagName   string    `json:"tag_name"`
	Color     *string   `json:"color,omitempty"`
	Icon      *string   `json:"icon,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
