package category

import (
	"context"
	"errors"
	"github.com/nogavadu/articles-service/internal/domain/converter"
	"github.com/nogavadu/articles-service/internal/domain/model"
	"github.com/nogavadu/articles-service/internal/repository"
	categoryRepo "github.com/nogavadu/articles-service/internal/repository/category"
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
		if errors.Is(err, categoryRepo.ErrInvalidArguments) {
			return 0, ErrInvalidArguments
		}
		if errors.Is(err, categoryRepo.ErrAlreadyExists) {
			return 0, ErrAlreadyExists
		}

		log.Error("failed to create category", slog.String("error", err.Error()))
		return 0, ErrInternalServerError
	}

	return id, nil
}

func (s *categoryService) GetAll(ctx context.Context, params *model.CategoryGetAllParams) ([]model.Category, error) {
	const op = "category.GetAll"
	log := s.log.With(slog.String("op", op))

	repoCategories, err := s.categoryRepo.GetAll(ctx, params)
	if err != nil {
		log.Error("failed to get categories", slog.String("error", err.Error()))
		return nil, ErrInternalServerError
	}

	categories := make([]model.Category, 0, len(repoCategories))
	for _, c := range repoCategories {
		categories = append(categories, *converter.ToCategory(&c))
	}

	return categories, nil
}

func (s *categoryService) GetById(ctx context.Context, id int) (*model.Category, error) {
	const op = "category.GetById"
	log := s.log.With(slog.String("op", op))

	repoCategory, err := s.categoryRepo.GetById(ctx, id)
	if err != nil {
		log.Error("failed to get category", slog.String("error", err.Error()))
		if errors.Is(err, categoryRepo.ErrNotFound) {
			return nil, ErrNotFound
		}
		if errors.Is(err, categoryRepo.ErrInvalidArguments) {
			return nil, ErrInvalidArguments
		}

		return nil, ErrInternalServerError
	}

	return converter.ToCategory(repoCategory), nil
}

func (s *categoryService) Update(ctx context.Context, id int, input *model.UpdateCategoryInput) error {
	const op = "category.Update"
	log := s.log.With(slog.String("op", op))

	if err := s.categoryRepo.Update(ctx, id, converter.ToRepoCategoryUpdateInput(input)); err != nil {
		log.Error("failed to update category", slog.String("error", err.Error()))
		return ErrInternalServerError
	}

	return nil
}
