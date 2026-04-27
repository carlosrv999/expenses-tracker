package model

import "time"

type Expense struct {
	ExpenseID       int64      `json:"expense_id"`
	CategoryID      int64      `json:"category_id"`
	PaymentMethodID int64      `json:"payment_method_id"`
	Currency        string     `json:"currency"`
	Amount          int64      `json:"amount"`
	ExpenseDate     time.Time  `json:"expense_date"`
	MerchantName    *string    `json:"merchant_name,omitempty"`
	Description     *string    `json:"description,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
	Tags            []Tag      `json:"tags,omitempty"`
}
