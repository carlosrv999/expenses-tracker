package expense

import "time"

type createRequest struct {
	CategoryID      int64     `json:"category_id"`
	PaymentMethodID int64     `json:"payment_method_id"`
	Currency        string    `json:"currency"`
	Amount          int64     `json:"amount"`
	ExpenseDate     time.Time `json:"expense_date"`
	MerchantName    *string   `json:"merchant_name"`
	Description     *string   `json:"description"`
	TagIDs          []int64   `json:"tag_ids"`
}

type updateRequest struct {
	CategoryID      int64     `json:"category_id"`
	PaymentMethodID int64     `json:"payment_method_id"`
	Currency        string    `json:"currency"`
	Amount          int64     `json:"amount"`
	ExpenseDate     time.Time `json:"expense_date"`
	MerchantName    *string   `json:"merchant_name"`
	Description     *string   `json:"description"`
	TagIDs          *[]int64  `json:"tag_ids"`
}
