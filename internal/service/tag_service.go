package service

import (
	"context"

	"expenses/internal/model"
	"expenses/internal/repository"
)

type TagService struct {
	repo *repository.TagRepository
}

func NewTagService(repo *repository.TagRepository) *TagService {
	return &TagService{repo: repo}
}

func (s *TagService) Create(ctx context.Context, t *model.Tag) error {
	return s.repo.Create(ctx, t)
}

func (s *TagService) Get(ctx context.Context, id int64) (*model.Tag, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TagService) List(ctx context.Context) ([]model.Tag, error) {
	return s.repo.List(ctx)
}

func (s *TagService) Update(ctx context.Context, t *model.Tag) error {
	return s.repo.Update(ctx, t)
}

func (s *TagService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
