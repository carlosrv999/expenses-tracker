package model

import "time"

type PaymentMethodType string

const (
	PaymentMethodCreditCard   PaymentMethodType = "credit_card"
	PaymentMethodCash         PaymentMethodType = "cash"
	PaymentMethodDebitCard    PaymentMethodType = "debit_card"
	PaymentMethodYape         PaymentMethodType = "yape"
	PaymentMethodPlin         PaymentMethodType = "plin"
	PaymentMethodBankTransfer PaymentMethodType = "bank_transfer"
)

func (t PaymentMethodType) Valid() bool {
	switch t {
	case PaymentMethodCreditCard, PaymentMethodCash, PaymentMethodDebitCard,
		PaymentMethodYape, PaymentMethodPlin, PaymentMethodBankTransfer:
		return true
	}
	return false
}

type PaymentMethod struct {
	PaymentMethodID int64             `json:"payment_method_id"`
	MethodName      string            `json:"method_name"`
	MethodType      PaymentMethodType `json:"method_type"`
	Icon            *string           `json:"icon,omitempty"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}
