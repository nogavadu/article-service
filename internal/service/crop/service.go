package crop

import (
	"context"
	"errors"
	"fmt"
	authService "github.com/nogavadu/articles-service/internal/clients/auth-service/grpc"
	"github.com/nogavadu/articles-service/internal/domain/converter"
	"github.com/nogavadu/articles-service/internal/domain/model"
	"github.com/nogavadu/articles-service/internal/repository"
	"github.com/nogavadu/articles-service/internal/repository/crop"
	"github.com/nogavadu/articles-service/internal/service"
	"github.com/nogavadu/platform_common/pkg/db"
	"log/slog"
)

var (
	ErrNotFound            = errors.New("crop not found")
	ErrAlreadyExists       = errors.New("crop already exists")
	ErrInvalidArguments    = errors.New("invalid article arguments")
	ErrInternalServerError = errors.New("internal server error")
	ErrAccessDenied        = errors.New("access denied")
)

type cropService struct {
	log *slog.Logger

	cropRepo           repository.CropRepository
	cropCategoriesRepo repository.CropCategoriesRepository
	statusRepo         repository.StatusRepository
	txManager          db.TxManager

	accessClient *authService.AccessServiceClient
	authClient   *authService.AuthServiceClient
}

func New(
	log *slog.Logger,
	cropRepository repository.CropRepository,
	cropCategoriesRepo repository.CropCategoriesRepository,
	statusRepo repository.StatusRepository,
	txManager db.TxManager,
	accessClient *authService.AccessServiceClient,
	authClient *authService.AuthServiceClient,
) service.CropService {
	return &cropService{
		log:                log,
		cropRepo:           cropRepository,
		cropCategoriesRepo: cropCategoriesRepo,
		statusRepo:         statusRepo,
		txManager:          txManager,
		accessClient:       accessClient,
		authClient:         authClient,
	}
}

func (s *cropService) Create(ctx context.Context, cropInfo *model.CropInfo) (int, error) {
	const op = "cropService.Create"
	log := s.log.With(slog.String("op", op))

	token, err := s.authClient.AccessToken(ctx)
	if err != nil {
		log.Error("failed to get access token", slog.String("error", err.Error()))
		return 0, ErrAccessDenied
	}

	status, err := s.statusRepo.GetByStatus(ctx, cropInfo.Status)
	fmt.Println(status)
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
		return 0, ErrAccessDenied
	}

	cropID, err := s.cropRepo.Create(ctx, converter.ToRepoCropInfo(cropInfo, status.Id))
	if err != nil {
		log.Error("failed to create crop", slog.String("error", err.Error()))

		if errors.Is(err, crop.ErrAlreadyExists) {
			return 0, ErrAlreadyExists
		}

		return 0, ErrInternalServerError
	}

	return cropID, nil
}

func (s *cropService) GetAll(ctx context.Context, params *model.CropGetAllParams) ([]model.Crop, error) {
	const op = "cropService.GetAll"
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

	repoCrops, err := s.cropRepo.GetAll(ctx, statusId)
	if err != nil {
		log.Error("failed to get crops", slog.String("error", err.Error()))
		return nil, ErrInternalServerError
	}

	crops := make([]model.Crop, 0, len(repoCrops))
	for _, repoCrop := range repoCrops {
		crops = append(crops, *converter.ToCrop(&repoCrop))
	}

	return crops, nil
}

func (s *cropService) GetById(ctx context.Context, id int) (*model.Crop, error) {
	const op = "cropService.GetById"
	log := s.log.With(slog.String("op", op))

	repoCrop, err := s.cropRepo.GetById(ctx, id)
	if err != nil {
		log.Error("failed to get crop", slog.String("error", err.Error()))
		return nil, ErrInternalServerError
	}

	return converter.ToCrop(repoCrop), nil
}

func (s *cropService) Update(ctx context.Context, id int, input *model.UpdateCropInput) error {
	const op = "cropService.Update"
	log := s.log.With(slog.String("op", op))

	token, err := s.authClient.AccessToken(ctx)
	if err != nil {
		log.Error("failed to get access token", slog.String("error", err.Error()))
		return ErrAccessDenied
	}
	err = s.accessClient.Check(ctx, token, authService.ModeratorAccessLevel)
	if err != nil {
		log.Error("access check failed", slog.String("error", err.Error()))
		return ErrAccessDenied
	}

	if err = s.cropRepo.Update(ctx, id, converter.ToRepoCropUpdateInput(input)); err != nil {
		log.Error("failed to update crop", slog.String("error", err.Error()))
		return ErrInternalServerError
	}

	return nil
}

func (s *cropService) Delete(ctx context.Context, id int) error {
	const op = "cropService.Delete"
	log := s.log.With(slog.String("op", op))

	token, err := s.authClient.AccessToken(ctx)
	if err != nil {
		log.Error("failed to get access token", slog.String("error", err.Error()))
		return ErrAccessDenied
	}
	err = s.accessClient.Check(ctx, token, authService.ModeratorAccessLevel)
	if err != nil {
		log.Error("access check failed", slog.String("error", err.Error()))
		return ErrAccessDenied
	}

	if err := s.cropRepo.Delete(ctx, id); err != nil {
		log.Error("failed to delete crop", slog.String("error", err.Error()))
		return ErrInternalServerError
	}

	return nil
}

func (s *cropService) AddRelation(ctx context.Context, cropId int, categoryId int) error {
	const op = "cropService.AddRelation"
	log := s.log.With(slog.String("op", op))

	token, err := s.authClient.AccessToken(ctx)
	if err != nil {
		log.Error("failed to get access token", slog.String("error", err.Error()))
		return ErrAccessDenied
	}
	err = s.accessClient.Check(ctx, token, authService.ModeratorAccessLevel)
	if err != nil {
		log.Error("access check failed", slog.String("error", err.Error()))
		return ErrAccessDenied
	}

	if err := s.cropCategoriesRepo.Create(ctx, cropId, categoryId); err != nil {
		log.Error("failed to add crop category", slog.String("error", err.Error()))
		return ErrInternalServerError
	}

	return nil
}

func (s *cropService) RemoveRelation(ctx context.Context, cropId int, categoryId int) error {
	const op = "cropService.RemoveRelation"
	log := s.log.With(slog.String("op", op))

	token, err := s.authClient.AccessToken(ctx)
	if err != nil {
		log.Error("failed to get access token", slog.String("error", err.Error()))
		return ErrAccessDenied
	}
	err = s.accessClient.Check(ctx, token, authService.ModeratorAccessLevel)
	if err != nil {
		log.Error("access check failed", slog.String("error", err.Error()))
		return ErrAccessDenied
	}

	if err := s.cropCategoriesRepo.Delete(ctx, cropId, categoryId); err != nil {
		log.Error("failed to remove crop category", slog.String("error", err.Error()))
		return ErrInternalServerError
	}

	return nil
}
