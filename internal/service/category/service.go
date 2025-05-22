package category

import (
	"context"
	"errors"
	"github.com/nogavadu/articles-service/internal/domain/converter"
	"github.com/nogavadu/articles-service/internal/domain/model"
	"github.com/nogavadu/articles-service/internal/repository"
	categoryRepo "github.com/nogavadu/articles-service/internal/repository/category"
	"github.com/nogavadu/articles-service/internal/service"
	"github.com/nogavadu/platform_common/pkg/db"
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

	categoryRepo       repository.CategoryRepository
	cropCategoriesRepo repository.CropCategoriesRepository
	txManager          db.TxManager
}

func New(
	log *slog.Logger,
	categoryRepo repository.CategoryRepository,
	cropCategoriesRepo repository.CropCategoriesRepository,
	txManager db.TxManager,
) service.CategoryService {
	return &categoryService{
		log:                log,
		categoryRepo:       categoryRepo,
		cropCategoriesRepo: cropCategoriesRepo,
		txManager:          txManager,
	}
}

func (s *categoryService) Create(
	ctx context.Context,
	categoryInfo *model.CategoryInfo,
	params *model.CategoryCreateParams,
) (int, error) {
	const op = "category.Create"
	log := s.log.With(slog.String("op", op))

	var id int
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		defer func() {
			if errTx != nil {
				log.Error("failed to create category", slog.String("error", errTx.Error()))
			}
		}()

		id, errTx = s.categoryRepo.Create(ctx, converter.ToRepoCategoryInfo(categoryInfo))
		if errTx != nil {
			if errors.Is(errTx, categoryRepo.ErrInvalidArguments) {
				return ErrInvalidArguments
			}
			if errors.Is(errTx, categoryRepo.ErrAlreadyExists) {
				return ErrAlreadyExists
			}

			return ErrInternalServerError
		}

		if params.CropId != nil {
			errTx = s.cropCategoriesRepo.Create(ctx, *params.CropId, id)
			if errTx != nil {
				return ErrInternalServerError
			}
		}

		return nil
	})

	return id, err
}

func (s *categoryService) GetAll(ctx context.Context, params *model.CategoryGetAllParams) ([]model.Category, error) {
	const op = "category.GetAll"
	log := s.log.With(slog.String("op", op))

	repoCategories, err := s.categoryRepo.GetAll(ctx, converter.ToRepoCategoryGetAllParams(params))
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

func (s *categoryService) Delete(ctx context.Context, id int) error {
	const op = "category.Delete"
	log := s.log.With(slog.String("op", op))

	if err := s.categoryRepo.Delete(ctx, id); err != nil {
		log.Error("failed to delete category", slog.String("error", err.Error()))
		return ErrInternalServerError
	}

	return nil
}
