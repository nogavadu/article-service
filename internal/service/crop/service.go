package crop

import (
	"context"
	"errors"
	"fmt"
	authService "github.com/nogavadu/articles-service/internal/clients/auth-service/grpc"
	"github.com/nogavadu/articles-service/internal/domain/converter"
	"github.com/nogavadu/articles-service/internal/domain/model"
	"github.com/nogavadu/articles-service/internal/repository"
	cropRepo "github.com/nogavadu/articles-service/internal/repository/crop"
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
	userClient   *authService.UserServiceClient
}

func New(
	log *slog.Logger,
	cropRepository repository.CropRepository,
	cropCategoriesRepo repository.CropCategoriesRepository,
	statusRepo repository.StatusRepository,
	txManager db.TxManager,
	accessClient *authService.AccessServiceClient,
	authClient *authService.AuthServiceClient,
	userClient *authService.UserServiceClient,
) service.CropService {
	return &cropService{
		log:                log,
		cropRepo:           cropRepository,
		cropCategoriesRepo: cropCategoriesRepo,
		statusRepo:         statusRepo,
		txManager:          txManager,
		accessClient:       accessClient,
		authClient:         authClient,
		userClient:         userClient,
	}
}

func (s *cropService) Create(ctx context.Context, userId int, cropInfo *model.CropInfo) (int, error) {
	const op = "cropService.Create"
	log := s.log.With(slog.String("op", op))

	token, err := s.authClient.AccessToken(ctx)
	if err != nil {
		log.Error("failed to get access token", slog.String("error", err.Error()))
		return 0, ErrAccessDenied
	}

	status, err := s.statusRepo.GetByStatus(ctx, cropInfo.Status)
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

	cropID, err := s.cropRepo.Create(ctx, converter.ToRepoCropInfo(cropInfo, status.Id, userId))
	if err != nil {
		log.Error("failed to create crop", slog.String("error", err.Error()))

		if errors.Is(err, cropRepo.ErrAlreadyExists) {
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
		statusId = 1
	}

	repoCrops, err := s.cropRepo.GetAll(ctx, statusId)
	if err != nil {
		log.Error("failed to get crops", slog.String("error", err.Error()))
		return nil, ErrInternalServerError
	}

	crops := make([]model.Crop, 0, len(repoCrops))
	for _, repoCrop := range repoCrops {
		repoStatus, err := s.statusRepo.GetById(ctx, repoCrop.Status)
		if err != nil {
			log.Error("failed to get status", slog.String("error", err.Error()))
			continue
		}

		crops = append(crops, *converter.ToCrop(&repoCrop, repoStatus.Status, nil))
	}

	return crops, nil
}

func (s *cropService) GetById(ctx context.Context, id int) (*model.Crop, error) {
	const op = "cropService.GetById"
	log := s.log.With(slog.String("op", op))

	var crop *model.Crop
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		defer func() {
			if errTx != nil {
				log.Error(fmt.Sprintf("failed to get crop by id: %v", errTx))
			}
		}()

		repoCrop, errTx := s.cropRepo.GetById(ctx, id)
		if errTx != nil {
			return ErrInternalServerError
		}

		repoStatus, errTx := s.statusRepo.GetById(ctx, repoCrop.Status)
		if errTx != nil {
			return ErrInternalServerError
		}

		var author *model.User
		if repoCrop.Author != nil {
			author, errTx = s.userClient.GetById(ctx, *repoCrop.Author)
			if errTx != nil {
				return ErrInternalServerError
			}
		}

		crop = converter.ToCrop(repoCrop, repoStatus.Status, author)

		return nil
	})

	return crop, err
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

	if input.Status != nil {
		status, _ := s.statusRepo.GetByStatus(ctx, *input.Status)
		if err = s.cropRepo.Update(ctx, id, converter.ToRepoCropUpdateInput(input, &status.Id)); err != nil {
			log.Error("failed to update crop", slog.String("error", err.Error()))
			return ErrInternalServerError
		}
	} else {
		if err = s.cropRepo.Update(ctx, id, converter.ToRepoCropUpdateInput(input, nil)); err != nil {
			log.Error("failed to update crop", slog.String("error", err.Error()))
			return ErrInternalServerError
		}
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
	err = s.accessClient.Check(ctx, token, authService.UserAccessLevel)
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
	err = s.accessClient.Check(ctx, token, authService.UserAccessLevel)
	if err != nil {
		log.Error("access check failed", slog.String("error", err.Error()))
		return ErrAccessDenied
	}

	if err = s.cropCategoriesRepo.Delete(ctx, cropId, categoryId); err != nil {
		log.Error("failed to remove crop category", slog.String("error", err.Error()))
		return ErrInternalServerError
	}

	return nil
}
