package service

import (
	"context"
	"errors"
	"fmt"

	"expenses/internal/model"
	"expenses/internal/repository"
)

var (
	ErrInvalidAmount   = errors.New("amount must be positive")
	ErrInvalidCurrency = errors.New("currency must be a 3-letter ISO 4217 code")
	ErrMismatchedLists = errors.New("internal error: expenses and tag lists length mismatch")
)

type ExpenseService struct {
	expenses *repository.ExpenseRepository
	tags     *repository.TagRepository
}

func NewExpenseService(expenses *repository.ExpenseRepository, tags *repository.TagRepository) *ExpenseService {
	return &ExpenseService{expenses: expenses, tags: tags}
}

func (s *ExpenseService) Create(ctx context.Context, e *model.Expense, tagIDs []int64) error {
	if err := validateExpense(e); err != nil {
		return err
	}
	if err := s.expenses.Create(ctx, e, tagIDs); err != nil {
		return err
	}
	return s.attachTags(ctx, e)
}

func (s *ExpenseService) Get(ctx context.Context, id int64) (*model.Expense, error) {
	e, err := s.expenses.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := s.attachTags(ctx, e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *ExpenseService) List(ctx context.Context, f repository.ExpenseFilter) ([]model.Expense, error) {
	expenses, err := s.expenses.List(ctx, f)
	if err != nil {
		return nil, err
	}
	// Tags are now populated by the repository (single query – no N+1)
	return expenses, nil
}

// ListPaginated returns a paginated list of expenses (with full metadata) and attaches
// the associated tags to every expense in the result — exactly like the non-paginated List method.
func (s *ExpenseService) ListPaginated(ctx context.Context, f repository.ExpenseFilter) (repository.PaginatedExpenseList, error) {
	result, err := s.expenses.ListPaginated(ctx, f)
	if err != nil {
		return repository.PaginatedExpenseList{}, err
	}
	// Tags are now populated by the repository (single query – no N+1)
	return result, nil
}

func (s *ExpenseService) Update(ctx context.Context, e *model.Expense, tagIDs *[]int64) error {
	if err := validateExpense(e); err != nil {
		return err
	}
	if err := s.expenses.Update(ctx, e, tagIDs); err != nil {
		return err
	}
	return s.attachTags(ctx, e)
}

func (s *ExpenseService) Delete(ctx context.Context, id int64) error {
	return s.expenses.SoftDelete(ctx, id)
}

func (s *ExpenseService) attachTags(ctx context.Context, e *model.Expense) error {
	tags, err := s.tags.ListByExpenseID(ctx, e.ExpenseID)
	if err != nil {
		return err
	}
	e.Tags = tags
	return nil
}

func validateExpense(e *model.Expense) error {
	if e.Amount <= 0 {
		return ErrInvalidAmount
	}
	if len(e.Currency) != 3 {
		return ErrInvalidCurrency
	}
	return nil
}

func (s *ExpenseService) BulkCreate(ctx context.Context, expenses []*model.Expense, tagIDsList [][]int64) error {
	if len(expenses) != len(tagIDsList) {
		return ErrMismatchedLists
	}

	for i, e := range expenses {
		if err := validateExpense(e); err != nil {
			return fmt.Errorf("row %d: %w", i, err)
		}
	}

	return s.expenses.BulkCreate(ctx, expenses, tagIDsList)
}
