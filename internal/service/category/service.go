package category

import (
	"context"
	"errors"
	authService "github.com/nogavadu/articles-service/internal/clients/auth-service/grpc"
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
	ErrAccessDenied        = errors.New("access denied")
)

type categoryService struct {
	log *slog.Logger

	categoryRepo       repository.CategoryRepository
	cropCategoriesRepo repository.CropCategoriesRepository
	statusRepo         repository.StatusRepository
	txManager          db.TxManager

	accessClient *authService.AccessServiceClient
	authClient   *authService.AuthServiceClient
	userClient   *authService.UserServiceClient
}

func New(
	log *slog.Logger,
	categoryRepo repository.CategoryRepository,
	cropCategoriesRepo repository.CropCategoriesRepository,
	statusRepo repository.StatusRepository,
	txManager db.TxManager,
	accessClient *authService.AccessServiceClient,
	authClient *authService.AuthServiceClient,
	userClient *authService.UserServiceClient,
) service.CategoryService {
	return &categoryService{
		log:                log,
		categoryRepo:       categoryRepo,
		cropCategoriesRepo: cropCategoriesRepo,
		statusRepo:         statusRepo,
		txManager:          txManager,
		accessClient:       accessClient,
		authClient:         authClient,
		userClient:         userClient,
	}
}

func (s *categoryService) Create(
	ctx context.Context,
	userId int,
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

		token, err := s.authClient.AccessToken(ctx)
		if err != nil {
			log.Error("failed to get access token", slog.String("error", err.Error()))
			return ErrAccessDenied
		}

		status, err := s.statusRepo.GetByStatus(ctx, categoryInfo.Status)
		if err != nil {
			log.Error("failed to get status", slog.String("error", err.Error()))
		}

		accessLevel := authService.ModeratorAccessLevel
		if status != nil && status.Id == 2 {
			accessLevel = authService.UserAccessLevel
		}

		err = s.accessClient.Check(ctx, token, accessLevel)
		if err != nil {
			log.Error("access check failed", slog.String("error", err.Error()))
			return ErrAccessDenied
		}

		id, errTx = s.categoryRepo.Create(ctx, converter.ToRepoCategoryInfo(categoryInfo, status.Id, userId))
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

	var statusId int
	if params.Status != nil {
		status, err := s.statusRepo.GetByStatus(ctx, *params.Status)
		if err != nil {
			log.Error("failed to get status", slog.String("error", err.Error()))
		}
		if status != nil {
			statusId = status.Id
		}
	} else {
		statusId = 2
	}

	repoCategories, err := s.categoryRepo.GetAll(ctx, converter.ToRepoCategoryGetAllParams(params, statusId))
	if err != nil {
		log.Error("failed to get categories", slog.String("error", err.Error()))
		return nil, ErrInternalServerError
	}

	categories := make([]model.Category, 0, len(repoCategories))
	for _, c := range repoCategories {
		repoStatus, err := s.statusRepo.GetById(ctx, c.Status)
		if err != nil {
			log.Error("failed to get status", slog.String("error", err.Error()))
			continue
		}

		var author *model.User
		if c.Author != nil {
			user, errTx := s.userClient.GetById(ctx, *c.Author)
			if errTx != nil {
				continue
			}
			author = user
		}

		categories = append(categories, *converter.ToCategory(&c, repoStatus.Status, author))
	}

	return categories, nil
}

func (s *categoryService) GetById(ctx context.Context, id int) (*model.Category, error) {
	const op = "category.GetById"
	log := s.log.With(slog.String("op", op))

	var category *model.Category
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		defer func() {
			if errTx != nil {
				log.Error("failed to get category", slog.String("error", errTx.Error()))
			}
		}()

		repoCategory, errTx := s.categoryRepo.GetById(ctx, id)
		if errTx != nil {
			log.Error("failed to get category", slog.String("error", errTx.Error()))
			if errors.Is(errTx, categoryRepo.ErrNotFound) {
				return ErrNotFound
			}
			if errors.Is(errTx, categoryRepo.ErrInvalidArguments) {
				return ErrInvalidArguments
			}

			return ErrInternalServerError
		}

		repoStatus, errTx := s.statusRepo.GetById(ctx, repoCategory.Status)
		if errTx != nil {
			return ErrNotFound
		}

		var author *model.User
		if repoCategory.Author != nil {
			user, errTx := s.userClient.GetById(ctx, *repoCategory.Author)
			if errTx != nil {
				return ErrNotFound
			}
			author = user
		}

		category = converter.ToCategory(repoCategory, repoStatus.Status, author)

		return nil
	})

	return category, err
}

func (s *categoryService) Update(ctx context.Context, id int, input *model.UpdateCategoryInput) error {
	const op = "category.Update"
	log := s.log.With(slog.String("op", op))

	if input.Status != nil {
		status, _ := s.statusRepo.GetByStatus(ctx, *input.Status)
		if err := s.categoryRepo.Update(ctx, id, converter.ToRepoCategoryUpdateInput(input, &status.Id)); err != nil {
			log.Error("failed to update category", slog.String("error", err.Error()))
			return ErrInternalServerError
		}
	} else {
		if err := s.categoryRepo.Update(ctx, id, converter.ToRepoCategoryUpdateInput(input, nil)); err != nil {
			log.Error("failed to update category", slog.String("error", err.Error()))
			return ErrInternalServerError
		}
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
