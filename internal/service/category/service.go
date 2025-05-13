package category

import (
	"context"
	"errors"
	"github.com/nogavadu/articles-service/internal/domain/converter"
	"github.com/nogavadu/articles-service/internal/domain/model"
	"github.com/nogavadu/articles-service/internal/repository"
	"github.com/nogavadu/articles-service/internal/service"
	"log/slog"
)

var (
	ErrNotFound            = errors.New("category not found")
	ErrAlreadyExists       = errors.New("category already exists")
	ErrInvalidArguments    = errors.New("invalid article arguments")
	ErrInternalServerError = errors.New("internal server error")
)

type categoryService struct {
	log *slog.Logger

	categoryRepo repository.CategoryRepository
}

func New(log *slog.Logger, categoryRepo repository.CategoryRepository) service.CategoryService {
	return &categoryService{
		log:          log,
		categoryRepo: categoryRepo,
	}
}

func (s *categoryService) Create(ctx context.Context, categoryInfo *model.CategoryInfo) (int, error) {
	const op = "category.Create"
	log := s.log.With(slog.String("op", op))

	id, err := s.categoryRepo.Create(ctx, converter.ToRepoCategoryInfo(categoryInfo))
	if err != nil {
		// TODO: add repo errors interceptors
		log.Error("failed to create category", slog.String("error", err.Error()))
		return 0, ErrInternalServerError
	}

	return id, nil
}

func (s *categoryService) GetList(ctx context.Context, cropId int) ([]*model.Category, error) {
	const op = "category.GetList"
	log := s.log.With(slog.String("op", op))

	repoCategories, err := s.categoryRepo.GetList(ctx, cropId)
	if err != nil {
		// TODO: add repo errors interceptors
		log.Error("failed to get categories", slog.String("error", err.Error()))
	}

	categories := make([]*model.Category, 0, len(repoCategories))
	for _, c := range repoCategories {
		categories = append(categories, converter.ToCategory(c))
	}

	return categories, nil
}

func (s *categoryService) GetAll(ctx context.Context) ([]*model.Category, error) {
	const op = "category.GetAll"
	log := s.log.With(slog.String("op", op))

	repoCategories, err := s.categoryRepo.GetAll(ctx)
	if err != nil {
		// TODO: add repo errors interceptors
		log.Error("failed to get categories", slog.String("error", err.Error()))
	}

	categories := make([]*model.Category, 0, len(repoCategories))
	for _, c := range repoCategories {
		categories = append(categories, converter.ToCategory(c))
	}

	return categories, nil
}
