package crop

import (
	"context"
	"errors"
	"fmt"
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
)

type cropService struct {
	log *slog.Logger

	cropRepo           repository.CropRepository
	cropCategoriesRepo repository.CropCategoriesRepository
	txManager          db.TxManager
}

func New(
	log *slog.Logger,
	cropRepository repository.CropRepository,
	cropCategoriesRepo repository.CropCategoriesRepository,
	txManager db.TxManager,
) service.CropService {
	return &cropService{
		log:                log,
		cropRepo:           cropRepository,
		cropCategoriesRepo: cropCategoriesRepo,
		txManager:          txManager,
	}
}

func (s *cropService) Create(ctx context.Context, cropInfo *model.CropInfo) (int, error) {
	const op = "cropService.Create"
	log := s.log.With(slog.String("op", op))

	fmt.Printf("SERVICE CROP INFO: %s\n", cropInfo)
	cropID, err := s.cropRepo.Create(ctx, converter.ToRepoCropInfo(cropInfo))
	if err != nil {
		log.Error("failed to create crop", slog.String("error", err.Error()))

		if errors.Is(err, crop.ErrAlreadyExists) {
			return 0, ErrAlreadyExists
		}

		return 0, ErrInternalServerError
	}

	return cropID, nil
}

func (s *cropService) GetAll(ctx context.Context) ([]model.Crop, error) {
	const op = "cropService.GetAll"
	log := s.log.With(slog.String("op", op))

	repoCrops, err := s.cropRepo.GetAll(ctx)
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

	if err := s.cropRepo.Update(ctx, id, converter.ToRepoCropUpdateInput(input)); err != nil {
		log.Error("failed to update crop", slog.String("error", err.Error()))
		return ErrInternalServerError
	}

	return nil
}

func (s *cropService) Delete(ctx context.Context, id int) error {
	const op = "cropService.Delete"
	log := s.log.With(slog.String("op", op))

	if err := s.cropRepo.Delete(ctx, id); err != nil {
		log.Error("failed to delete crop", slog.String("error", err.Error()))
		return ErrInternalServerError
	}

	return nil
}

func (s *cropService) AddRelation(ctx context.Context, cropId int, categoryId int) error {
	const op = "cropService.AddRelation"
	log := s.log.With(slog.String("op", op))

	if err := s.cropCategoriesRepo.Create(ctx, cropId, categoryId); err != nil {
		log.Error("failed to add crop category", slog.String("error", err.Error()))
		return ErrInternalServerError
	}

	return nil
}

func (s *cropService) RemoveRelation(ctx context.Context, cropId int, categoryId int) error {
	const op = "cropService.RemoveRelation"
	log := s.log.With(slog.String("op", op))

	if err := s.cropCategoriesRepo.Delete(ctx, cropId, categoryId); err != nil {
		log.Error("failed to remove crop category", slog.String("error", err.Error()))
		return ErrInternalServerError
	}

	return nil
}
