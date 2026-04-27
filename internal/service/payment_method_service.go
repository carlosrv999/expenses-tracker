package service

import (
	"context"
	"errors"

	"expenses/internal/model"
	"expenses/internal/repository"
)

var ErrInvalidPaymentMethodType = errors.New("invalid payment method type")

type PaymentMethodService struct {
	repo *repository.PaymentMethodRepository
}

func NewPaymentMethodService(repo *repository.PaymentMethodRepository) *PaymentMethodService {
	return &PaymentMethodService{repo: repo}
}

func (s *PaymentMethodService) Create(ctx context.Context, p *model.PaymentMethod) error {
	if !p.MethodType.Valid() {
		return ErrInvalidPaymentMethodType
	}
	return s.repo.Create(ctx, p)
}

func (s *PaymentMethodService) Get(ctx context.Context, id int64) (*model.PaymentMethod, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *PaymentMethodService) List(ctx context.Context) ([]model.PaymentMethod, error) {
	return s.repo.List(ctx)
}

func (s *PaymentMethodService) Update(ctx context.Context, p *model.PaymentMethod) error {
	if !p.MethodType.Valid() {
		return ErrInvalidPaymentMethodType
	}
	return s.repo.Update(ctx, p)
}

func (s *PaymentMethodService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
